<template>

</template>

<script>
 import {ref, onMounted} from "vue";

 export default {
  name: "4.组合式API"
}

const count = ref(0)

 function increment() {
   count.value++
 }

 onMounted(()=>{
   console.log('The initial count is ${count.value}.')
     }
 )
</script>

<style scoped>

</style>