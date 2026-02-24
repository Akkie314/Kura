<script setup lang="ts">
const emit = defineEmits<{
    (e: 'change', file: File): void
}>();

defineProps<{
    href?: string;
    disabled?: boolean;
    download?: boolean | string;
}>();

function onChange(e: Event) {
    const target = e.target as HTMLInputElement;
    const files = target.files;

    if (files && files.length > 0) {
        emit('change', files[0]);
        target.value = "";
    }
} 
</script>

<template>
    <div>
        <a v-if="href" :href="disabled ? undefined : href"
            :download="download === false ? undefined : (download === true ? '' : download)" class="btn-base"
            :class="{ 'is-disabled': disabled }">
            <slot></slot>
        </a>

        <label v-else class="btn-base" :class="{ 'is-disabled': disabled }" @click.stop>
            <input type="file" @change="onChange" :disabled="disabled" name="selectedFile" tabindex="0"
                class="visually-hidden">
            <slot></slot>
        </label>
    </div>
</template>

<style scoped>
.btn-base {
    display: flex;
    padding: 0.75rem 1.5rem;
    justify-content: center;
    align-items: center;
    gap: 10px;

    border-radius: 0.5rem;
    background: var(--color-bg-blue);
    box-shadow: var(--box-shadow-value);

    cursor: pointer;
    transition: background 0.15s ease-out;

    text-decoration: none;
    color: inherit;
    box-sizing: border-box;
}

.btn-base:hover,
.btn-base:active,
.btn-base:focus,
.btn-base:has(input:focus) {
    background: var(--color-bg-dark-blue);
}

.btn-base.is-disabled {
    opacity: 0.5;
    cursor: not-allowed;
    pointer-events: none;
}

.visually-hidden {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    white-space: nowrap;
    border-width: 0;
}
</style>