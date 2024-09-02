<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script>
	import { debounce} from "lodash";
  export default {
  name: "App",
  mounted() {
    // 保存原始的 ResizeObserver 类
    const OriginalResizeObserver = window.ResizeObserver;
    // 重写 ResizeObserver 类
    window.ResizeObserver = class DebouncedResizeObserver extends OriginalResizeObserver {
      constructor(callback) {
        // 对回调函数进行防抖处理，延迟 100 毫秒执行
        const debouncedCallback = debounce(callback, 100);
        // 调用父类的构造函数，并传入经过防抖处理的回调函数
        super(debouncedCallback);
      }
    };
  }
};
</script>

<style>

</style>
