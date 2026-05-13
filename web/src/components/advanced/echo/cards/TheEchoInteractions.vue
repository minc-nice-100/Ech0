<!-- SPDX-License-Identifier: AGPL-3.0-or-later -->
<!-- Copyright (C) 2025-2026 lin-snow -->
<template>
  <section class="echo-interactions w-full max-w-sm mx-auto">
    <TheComment />
  </section>
</template>

<script setup lang="ts">
import TheComment from '@/components/advanced/TheComment.vue'
</script>

<style scoped>
.echo-interactions {
  background-image: linear-gradient(
    to right,
    var(--color-border-subtle) 0,
    var(--color-border-subtle) 5px,
    transparent 5px,
    transparent 8px
  );
  background-size: 8px 1px;
  background-repeat: repeat-x;
  background-position: left top;
}
</style>
