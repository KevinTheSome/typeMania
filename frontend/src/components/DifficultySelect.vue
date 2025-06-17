<script setup lang="ts">
import { ref, defineEmits } from "vue";
const emit = defineEmits(["difficulty-event"]);
function diffHandeler(evt: Event) {
  const target = evt.target as HTMLSelectElement;
  emit("difficulty-event", target.value);
}
</script>

<template>
  <select
    name="difficulty"
    id=""
    @change="diffHandeler"
    @keydown.prevent
    class="bg-transparent border-none outline-none text-aqua-500 font-semibold"
  >
    <option disabled selected class="text-gray-500">Select Difficulty</option>
    <option value="2" class="text-white">Easy</option>
    <option value="100" class="text-white">Normal</option>
    <option value="150" class="text-white">Hard</option>
    <option value="300" class="text-white">HardCore</option>
  </select>
</template>
