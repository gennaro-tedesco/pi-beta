# Fix: sluggish screen transitions / animations on high-resolution + low-power displays (Raspberry Pi)

## Problem

On a Raspberry Pi, the app is sluggish and screen transitions are choppy at high screen
resolution, and becomes responsive again when the window is made smaller. This points to
CSS operations whose cost scales with pixel count (paint/rasterization), rather than JS
computation (the app has no canvas/JS rendering loop — all animation is CSS/DOM driven).

## Root causes

1. **`filter: brightness()` on a full-viewport element**
   `frontend/src/lib/views/WeatherView.svelte:449-451`
   ```css
   .weather-widget--night {
     filter: brightness(var(--weather-widget-night-brightness));
   }
   ```
   `.weather-widget` is `position: absolute; inset: 0` at 100% width/height
   (`WeatherView.svelte:437-442`), i.e. full screen when open. Any CSS `filter` forces the
   browser to rasterize the whole subtree into an offscreen bitmap every time it (or its
   animating children) repaint. Cost scales with viewport pixel area.

2. **Animated `background-position` on the full-screen root**
   `frontend/src/App.svelte:183-186`
   ```css
   .screen {
     background-image: var(--app-background);
     background-size: 200% 200%;
     animation: driftGradient 30s ease-in-out infinite;
   }
   ```
   `background-position` is a paint-triggering property (unlike `transform`/`opacity`), so
   this repaints the full-viewport gradient layer every frame, continuously, for the whole
   30s loop.

3. **Sequential (non-overlapping) full-screen `scale` transitions on every navigation**
   Same pattern repeated in every view root:
   - `frontend/src/App.svelte:108-109`
   - `frontend/src/lib/views/WeatherView.svelte:289-290`
   - `frontend/src/lib/views/PhotoSlideshow.svelte:87-88`
   - `frontend/src/lib/views/NetworkView.svelte:168-169`
   - `frontend/src/lib/views/CalendarView.svelte:260-261`
   - `frontend/src/lib/views/CityPicker.svelte:264-265`
   ```js
   out:scale={{ duration: transitionDuration }}
   in:scale={{ duration: transitionDuration, delay: transitionDuration }}
   ```
   `transitionDuration = 1200` (`frontend/src/lib/transition.ts:1`). The `delay` on `in`
   makes it wait for `out` to finish, so every navigation is 2.4s of continuous
   scale/opacity animation on a full-viewport element — and issues #1/#2 keep repainting
   the pixels underneath that animation the whole time it plays.

   `WeatherView.svelte` also drives some of its own transitions off a mutable
   `weatherViewTransitionDuration` (`WeatherView.svelte:76,93,210`), which is the same
   `out`/`in` + `delay` pattern and should be fixed the same way.

## Proposed changes (implement in this order)

### 1. Replace `filter: brightness()` with an opacity overlay

File: `frontend/src/lib/views/WeatherView.svelte`

Remove the `filter` from `.weather-widget--night`. Add an absolutely positioned overlay
element covering the widget (`position: absolute; inset: 0`) with a dark background color
and an opacity value that reproduces the same visual dimming as the current
`brightness(0.82)`, toggled by the existing `weather-widget--night` state. `opacity` is
compositor-only (no repaint), unlike `filter`.

### 2. Replace animated `background-position` with `transform`

File: `frontend/src/App.svelte`

Move the gradient onto its own absolutely-positioned child layer, oversized relative to its
container (so it can pan without exposing an edge), and animate that layer's `transform:
translate(...)` in the `driftGradient` keyframes instead of `background-position`. Keep the
visual drift equivalent to today's. `transform` is compositor-only.

### 3. Overlap `in`/`out` transitions instead of running them sequentially

Files: every `in:scale`/`out:scale` pair listed under root cause #3.

Remove the `delay: transitionDuration` argument from each `in:scale` call so the incoming
and outgoing views crossfade/scale simultaneously instead of sequentially. This is only
safe to do **after** steps 1 and 2: once both animating layers are purely
`transform`/`opacity` (no `filter`, no `background-position`), compositing two of them at
once is cheap (GPU blends two pre-rasterized textures). Doing this before 1/2 would make
things worse, since it would double up the still-repainting layers instead of just
composited ones.

### 4. (Optional, only if still not smooth after 1-3) Shorten `transitionDuration`

File: `frontend/src/lib/transition.ts:1`

`transitionDuration` is currently `1200`. If transitions are still not snappy enough on the
Pi after 1-3, consider lowering this value. Treat as a secondary tuning knob, not the
primary fix.

## Constraints (this repo's `AGENTS.md`, mandatory for every change)

- No code comments.
- No hardcoded literals — extract every literal (numbers, colors, durations) into a named
  constant / CSS custom property, consistent with the existing `--*` custom property
  pattern used throughout these files.
- Every component must remain scalable to container/viewport size, not fixed dimensions —
  verify at multiple widths/heights.
- Reuse existing components/patterns; do not introduce new UI components.
- Maintain spatial symmetry (equal opposing margins/padding) unless a reason is stated.
- Commit message: a single concise line, not a changelog.

## Verification

- Visual: night dimming and background drift must look the same as before the change.
- Performance: test on the actual Raspberry Pi at native (high) resolution — screen
  transitions and the weather widget's night mode should no longer drop frames or feel
  sluggish.
- Regression check: confirm behavior is unchanged at smaller window sizes (where it was
  already smooth).
- Re-verify `prefers-reduced-motion` handling still works (see existing media query in
  `App.svelte:503-511` and `WeatherView.svelte` reduced-motion rules).
