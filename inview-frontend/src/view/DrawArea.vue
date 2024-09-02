<template>
  <div>
    <div class="tool-box">
      <el-button :disabled="isDrawing" @click="startDraw">绘制区域</el-button>
      <el-button :disabled="isDrawing" @click="clearAll">清除</el-button>
    </div>
    <div class="canvas-wrap">
      <canvas id="imgCanvas" ref="canvaxbox"></canvas>
      <canvas id="drawCanvas" ref="canvas" :style="{ cursor: isDrawing ? 'crosshair' : 'default' }"> </canvas>
      <canvas id="saveCanvas" ref="canvasSave"></canvas>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      imgUrl: '',
      isDrawing: false, // 是否正在绘制
      ratio: 1,
      imgWidth: 0,
      imgHeight: 0,
      wrapWidth: 300,
      wrapHeight: 300,
      canvasWidth: 300,
      canvasHeight: 300,
      drawingPoints: [],
      drawnPoints: [],
      imgCanvas: null,
      imgCtx: null,
      drawCanvas: null,
      drawCtx: null,
      saveCanvas: null,
      saveCtx: null,
      submitData: []
    }
  },
  mounted() {
    this.initCanvas();
    this.getImage();
  },
  props: {
    screenshotData: {
      type: Object,
    },
    coordinates: { // 新增坐标 prop
      type: Array,
      default: () => []
    }
  },
  methods: {
    //绘制之前的坐标
    renderInitialCoordinates() {
      if (this.coordinates.length > 0) {
        this.drawnPoints = this.coordinates.map(area => {
          return area.map(point => {
            return [parseFloat(point.x / this.ratio * this.imgWidth), parseFloat(point.y / this.ratio * this.imgHeight)]; // 转换为画布坐标
          });
        });
        this.drawnPoints.forEach(points => {
          this.drawSaveArea(points); // 绘制区域
        });
      }
    },
    initCanvas() { // 初始化canvas画布
      let canvasWrap = document.getElementsByClassName('canvas-wrap');
      this.wrapWidth = canvasWrap[0].clientWidth;
      this.wrapHeight = canvasWrap[0].clientHeight;

      this.imgCanvas = document.getElementById('imgCanvas');
      this.imgCtx = this.imgCanvas.getContext('2d');

      // 绘制canvas
      this.drawCanvas = document.getElementById('drawCanvas');
      this.drawCtx = this.drawCanvas.getContext('2d');

      // 保存绘制区域 saveCanvas
      this.saveCanvas = document.getElementById('saveCanvas');
      this.saveCtx = this.saveCanvas.getContext('2d');
    },
    getImage() {
      this.imgUrl = this.screenshotData;
      if (this.imgUrl) {
        const image = new Image();
        image.src = this.imgUrl;
        image.onload = () => {
          this.imgWidth = image.width;
          this.imgHeight = image.height;
          this.initImgCanvas();
        };
      }
    },
    initImgCanvas() {
      // 计算宽高比
      let ww = this.wrapWidth; // 画布宽度
      let wh = this.wrapHeight; // 画布高度
      let iw = this.imgWidth; // 图片宽度
      let ih = this.imgHeight; // 图片高度

      if (iw / ih < ww / wh) { // 以高为主
        this.ratio = ih / wh;
        this.canvasHeight = wh;
        this.canvasWidth = wh * iw / ih;
      } else { // 以宽为主
        this.ratio = iw / ww;
        this.canvasWidth = ww;
        this.canvasHeight = ww * ih / iw;
      }

      // 初始化画布大小
      this.imgCanvas.width = this.canvasWidth;
      this.imgCanvas.height = this.canvasHeight;
      this.drawCanvas.width = this.canvasWidth;
      this.drawCanvas.height = this.canvasHeight;
      this.saveCanvas.width = this.canvasWidth;
      this.saveCanvas.height = this.canvasHeight;

      // 图片加载绘制
      let img = document.createElement('img');
      img.src = this.imgUrl;
      img.onload = () => {
        console.log('图片已加载');
        this.imgCtx.drawImage(img, 0, 0, this.canvasWidth, this.canvasHeight);
        this.renderInitialCoordinates() //渲染原始数据
      };

    },
    startDraw() { // 绘制区域
      if (this.isDrawing) return;
      this.isDrawing = true;
      // 绘制逻辑
      this.drawCanvas.addEventListener("click", this.drawImageClickFn);
      this.drawCanvas.addEventListener("mousemove", this.drawImageMoveFn);
    },
    clearAll() { // 清空所有绘制区域
      this.saveCtx.clearRect(0, 0, this.canvasWidth, this.canvasHeight);
      this.drawnPoints = [];
    },

    drawImageClickFn(e) {
      if (e.offsetX || e.layerX) {
        let pointX = e.offsetX == undefined ? e.layerX : e.offsetX;
        let pointY = e.offsetY == undefined ? e.layerY : e.offsetY;

        // 如果点击位置接近第一个点，则停止绘制
        if (this.drawingPoints.length > 2) {
          let firstPoint = this.drawingPoints[0];
          if (Math.abs(pointX - firstPoint[0]) < 10 && Math.abs(pointY - firstPoint[1]) < 10) {
            this.drawImageDblClickFn();
            return;
          }
        }

        let lastPoint = this.drawingPoints[this.drawingPoints.length - 1] || [];
        if (lastPoint[0] !== pointX || lastPoint[1] !== pointY) {
          this.drawingPoints.push([pointX, pointY]);
        }
      }
    },
    drawImageMoveFn(e) {
      let drawCtx = this.drawCtx;
      if (e.offsetX || e.layerX) {
        let pointX = e.offsetX == undefined ? e.layerX : e.offsetX;
        let pointY = e.offsetY == undefined ? e.layerY : e.offsetY;
        // 绘制
        drawCtx.clearRect(0, 0, this.canvasWidth, this.canvasHeight);

        // 绘制点
        drawCtx.fillStyle = 'blue';
        this.drawingPoints.forEach((item) => {
          drawCtx.beginPath();
          drawCtx.arc(...item, 6, 0, 180);
          drawCtx.fill(); //填充
        });

        // 绘制动态区域
        drawCtx.save();
        drawCtx.beginPath();
        this.drawingPoints.forEach((item) => {
          drawCtx.lineTo(...item);
        });
        drawCtx.lineTo(pointX, pointY);
        drawCtx.lineWidth = "3";
        drawCtx.strokeStyle = "blue";
        drawCtx.fillStyle = 'rgba(255, 0, 0, 0.3)';
        drawCtx.stroke();
        drawCtx.fill(); //填充
        drawCtx.restore();
      }
    },
    drawImageDblClickFn() {
      let drawCtx = this.drawCtx;
      // 清空绘制图层
      drawCtx.clearRect(0, 0, this.canvasWidth, this.canvasHeight);
      // 绘制区域至保存图层
      this.drawSaveArea(this.drawingPoints);

      this.drawnPoints.push([...this.drawingPoints]);
      this.drawingPoints = [];
      this.isDrawing = false;

      // 绘制结束逻辑
      this.drawCanvas.removeEventListener("click", this.drawImageClickFn);
      this.drawCanvas.removeEventListener("mousemove", this.drawImageMoveFn);
    },
    drawSaveArea(points) {
      if (points.length === 0) return;
      this.saveCtx.save();
      this.saveCtx.beginPath();
      points.forEach((item) => {
        this.saveCtx.lineTo(...item);
      });
      this.saveCtx.closePath();
      this.saveCtx.lineWidth = "2";
      this.saveCtx.fillStyle = 'rgba(255,0, 255, 0.3)';
      this.saveCtx.strokeStyle = "red";
      this.saveCtx.stroke();
      this.saveCtx.fill(); //填充
      this.saveCtx.restore();
    },
    savePoints() { // 将画布坐标数据转换成提交数据
      let objectPoints = [];
      this.drawnPoints.forEach(area => {
        let points = [];
        area.forEach((point) => {
          let polygon = {};
          polygon.x = parseFloat((point[0] * this.ratio / this.imgWidth));
          polygon.y = parseFloat((point[1] * this.ratio / this.imgHeight));
          points.push(polygon);
        });
        objectPoints.push(points);
      });
      this.submitData = objectPoints;
      return objectPoints;
    },
  }
};
</script>

<style scoped>
.tool-box {
  width: 100%;
  box-sizing: border-box;
  text-align: right;
}

.canvas-wrap {
  width: 80vw;
  height: 35vw;
  width: 60vw;
  height: 33.75vw;
  margin: 0px auto;
  background-color: #000;
  border: 3px;
  border-color: #333;
  position: relative;
}

#imgCanvas,
#drawCanvas,
#saveCanvas {
  background: rgba(255, 0, 255, 0);
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

#drawCanvas {
  z-index: 2;
}
</style>
