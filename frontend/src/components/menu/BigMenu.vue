<template>
  <div class="wrapper" :class="{ 'hvr-shutter-in-horizontal': !isActive, selected: isActive }"v-on:click="select($event)">
      <div class="label">{{ label }}</div>
      <div class="icon" ref="icon" :class="{ wiggle: isActive }"><icon :name="iconName" :w="30" :h="30"></icon></div>
  </div>
</template>

<script>
import icon from 'vue-icon'
export default {
  name: 'big-menu',
  components: { icon },
  props: ['navigation', 'label', 'iconName', 'currentModule'],
  data () {
    return {
      msg: '',
      isActive: false
    }
  },
  methods: {
    select: function (event) {
      this.$parent.updateCurrentModule(this.label)
      this.$router.push({ path: this.navigation })
    }
  },
  watch: {
    currentModule: function () {
      if (this.currentModule === this.label) {
        this.isActive = true
      } else {
        this.isActive = false
      }
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
div .wrapper {
  width: 150px;
  height: 74px;
  /*border-top: 0.5px solid #ececec;*/
  border-bottom: 0.5px solid #ececec;
  position: relative;
  cursor: pointer;
  background: #00B5AD;
}

.hvr-shutter-in-horizontal {
  display: block;
  background: #00B5AD;
}
.hvr-shutter-in-horizontal:before {
  background: #6ca893;
}

@keyframes wiggle {
/*  10%, 90% {
    transform: translate3d(-1px, 0, 0);
  }
  20%, 80% {
    transform: translate3d(2px, 0, 0);
  }
  30%, 50%, 70% {
    transform: translate3d(-4px, 0, 0);
  }
  40%, 60% {
    transform: translate3d(4px, 0, 0);
  }*/
  50% {
    transform: translate3d(-8px, 0, 0);
  }
}

div .wrapper:hover .icon{
  fill: white;
  display: inline-block;
  animation: wiggle 0.5s infinite;
}

.wiggle {
  display: inline-block;
  animation: wiggle 0.5s infinite;
}

.selected {
  color: white;
}
.selected .icon{
  fill: white;
}

.label {
  margin-left: 5px;
  position: absolute;
  top: 42%;
}

.icon {
  margin-right: 5px;
  position: absolute;
  top: 35%;
  right: 0;
  fill: #666;
  display:inline-block
}
</style>
