<script setup lang="ts">
defineProps<{
    value: number
}>()
</script>



<template>
    <div class="progress-container">
        <div class="progress-bar" :style="{width:value + '%'}"></div>
    </div>
</template>



<style scoped>
div.progress-container {
    width: 100%;
    max-width: 600px;
    height: 3rem;

    overflow: hidden;

    background-color: var(--color-bg-black);
    border: 1px solid var(--color-text-white);
    border-radius: 1rem;
}

div.progress-bar {
    width: 0;
    height: 100%;

    background-color: var(--color-bg-blue);

    transition: width 1.5s ease-out;
}
</style>