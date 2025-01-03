<template>
  <div>
    <input class="pixelated-input" type="text" v-bind="$attrs" v-model="text"></input>
  </div>
</template>

<script setup>
	import { defineModel } from 'vue'
	const text = defineModel({
		type: String,
		required: false,
		default: ''
	})
</script>

<style scoped>
.pixelated-input{
  border: 3px dashed black;
  padding: 3px;
  padding-left: 8px;
  border-radius: 8px;
  font-family: 'VT323';
  font-size: 25px;
}
.pixelated-input:focus {
    outline-width: 0;
    outline: none;
}

</style>
