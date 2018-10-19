<template lang="html">
<!-- <transition name="modal"> -->
   <div class="modal">
    <div class="modal-window">
      <p>{{modalTitle}}</p>
      <div class="actions">
       <button class="cancel" @click="cancel">Cancel</button>
       <button class="confirm" @click="confirm">Confirm</button>
      </div>
    </div>
   </div>
<!-- </transition> -->
</template>
<script>
export default {
  name: 'modal',
  props: ['modalTitle', 'modalFunction', 'record'],
  data () {
    return {}
  },
  methods: {
    modalShow (rec) {
      this.record = rec
    },
    confirm () {
      this.modalFunction(this.record)
      this.$parent.needToConfirm = false
    },
    cancel () {
      this.$parent.needToConfirm = false
    }
  }
}
</script>
<style scoped>
  button {
   cursor: pointer;
  }
  
  .modal {
   position: fixed;
   top: 0;
   left: 0;
   bottom: 0;
   right: 0;
   background-color: rgba(0, 0, 0, .5);
   /*opacity:1.0;*/
/*   -webkit-filter: blur(2px);
       -moz-filter: blur(2px);
       -o-filter: blur(2px);
       -ms-filter: blur(2px);
       filter: blur(2px);*/
  }
  
  .modal-window {
   position: absolute;
   top: 30%;
   left: 50%;
   transition: .5s;
   width: 100%;
   min-width: 400px;
   max-width: 600px;
   background: white;
   box-shadow: 0 0 10px rgba(0, 0, 0, .5);
   transform: translate(-50%, -50%);
   padding: 1em;
   color: black;
   text-align: center;
  }
  
  .modal-window .actions {
   /*display: flex;*/
   /*justify-content: flex-end;*/
  }
  
  .actions button {
   font-size: inherit;
   margin: 4px;
   border: none;
   padding: 6px 8px;
   cursor: pointer;
  }
  
  .actions .cancel {
   background: darkred;
   color: white;
  }
  
  .actions .confirm {
   background: darkcyan;
   color: white;
  }
</style>
