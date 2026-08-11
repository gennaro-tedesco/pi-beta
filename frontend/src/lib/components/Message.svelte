<script lang="ts">
  export let variant: 'error' | 'warning' | 'success' = 'success'
  export let message: string = ''
</script>

<div class="message message--{variant}" role="alert">
  <slot>{message}</slot>
</div>

<style>
  .message {
    padding: 0.75rem 1rem;
    border-radius: 6px;
    font-size: 0.95rem;
    text-align: left;
  }

  .message--error {
    background-color: rgba(220, 53, 69, 0.15);
    border: 1px solid rgba(220, 53, 69, 0.6);
    color: #ff8a94;
  }

  .message--warning {
    background-color: rgba(255, 193, 7, 0.15);
    border: 1px solid rgba(255, 193, 7, 0.6);
    color: #ffd764;
  }

  .message--success {
    background-color: rgba(40, 167, 69, 0.15);
    border: 1px solid rgba(40, 167, 69, 0.6);
    color: #7be495;
  }
</style>
