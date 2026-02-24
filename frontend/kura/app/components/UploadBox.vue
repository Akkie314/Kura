<script setup lang="ts">
import { ref } from 'vue';

const selectedFile = ref<File | null>(null);
const isDragging = ref<boolean>(false);
const dragCounter = ref<number>(0);
const fileInput = ref<HTMLElement | null>(null);

const emit = defineEmits<{
    (e: 'change', file: File): void
}>();

function handlerChangeFile(file:File) {
    selectedFile.value = file;
    emit('change', file);
}

function onDrop(e:DragEvent) {
    onDragLeave();

    const files = e.dataTransfer?.files;

    if (files && files.length > 0) {
        handlerChangeFile(files[0]);
    }
}

function onDragEnter() {
    dragCounter.value++;
    isDragging.value = true;
}

function onDragLeave() {
    dragCounter.value--;

    if (dragCounter.value === 0) {
        isDragging.value = false;
    }
}

// アップロードボックスをクリックした時のファイル選択用
function onClickBox() {
    fileInput.value?.click();
}

function onChange(e:Event) {
    const target = e.target as HTMLInputElement;
    const files = target.files;

    if (files && files.length > 0) {
        handlerChangeFile(files[0]);
        target.value = ""
    } 
}
</script>

<template>
    <div class="upload-box"
        :class="{'isDragging':isDragging}"
        @click="onClickBox"
        @dragenter="onDragEnter" 
        @dragleave="onDragLeave" 
        @drop.prevent="onDrop" 
        @dragover.prevent>
        <img src="/icon/Upload.svg" alt="" class="upload-icon">
        <p>ここにファイルをドロップ</p>
        <p class="gray-text">または</p>
        <AppButton type="file" @change="handlerChangeFile">ファイルを選択</AppButton>
    </div>

    <!-- アップロードボックスをクリックした時のファイル選択用 -->
    <input type="file" @change="onChange" ref="fileInput" class="hidden-file-input">
</template>

<style scoped>
div.upload-box {
    display: flex;
    width: 100%;
    max-width: 600px;
    height: 50%;

    padding: 2rem;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    gap: 1rem;

    border-radius: 1rem;
    border: 2px dashed var(--color-text-light-gray);
    background-color: var(--color-bg-black);

    transition: transform 0.15s ease-out;
}

div.upload-box.isDragging {
    transform: scale(1.1);
}

img.upload-icon {
    display: flex;
    width: 4rem;
    height: 4rem;

    padding: 1rem;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    flex-shrink: 0;

    border-radius: 100%;
    background: rgba(255, 255, 255, 0.10);
}

p.gray-text {
    color: var(--color-text-light-gray);
    font-size: 0.75rem;
}
</style>