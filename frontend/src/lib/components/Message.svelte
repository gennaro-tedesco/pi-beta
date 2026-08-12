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
    background-color: rgba(255, 205, 210, 0.6);
    border: 1px solid rgba(198, 40, 57, 0.35);
    color: #a3283a;
  }

  .message--warning {
    background-color: rgba(255, 236, 179, 0.65);
    border: 1px solid rgba(179, 130, 0, 0.35);
    color: #8a6100;
  }

  .message--success {
    background-color: rgba(200, 239, 214, 0.6);
    border: 1px solid rgba(43, 130, 82, 0.35);
    color: #2f7a52;
  }
</style>
