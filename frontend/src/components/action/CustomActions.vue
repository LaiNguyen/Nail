<template>
  <div class="custom-actions">
    <modal v-if="needToConfirm" :modalTitle="modalTitle" :modalFunction="modalFunction" :record="record"></modal>
<!--     <button class="ui basic button" @click="itemAction('view', rowData, rowIndex)">
      <span class="ui green label"><i class="large zoom icon" style="margin: 0px"></i></span>
    </button> -->
    <button class="ui basic button" @click="itemAction('edit', rowData, rowIndex)">
      <span class="ui orange label"><i class="large edit icon" style="margin: 0px"></i></span>
    </button>
    <button class="ui basic button" @click="itemAction('delete', rowData, rowIndex)">
      <span class="ui red label"><i class="large remove icon" style="margin: 0px"></i></span>
    </button>
  </div>
</template>
<script>
  import config from '@/config.js'
  import axios from 'axios'
  import modal from '@/mixins/modal.vue'
  export default {
    data () {
      return {
        needToConfirm: false,
        modalTitle: '',
        record: {},
        modalFunction: function () {

        }
      }
    },
    // mixins: [modal],
    components: { config, axios, modal },
    props: {
      rowData: {
        type: Object,
        required: true
      },
      rowIndex: {
        type: Number
      }
    },
    methods: {
      itemAction (action, data, index) {
        switch (action) {
          case 'view':
            this.$parent.$parent.onViewClicked(data)
            break
          case 'edit':
            this.$parent.$parent.onEditClicked(data)
            break
          case 'delete':
            this.needToConfirm = true
            this.modalTitle = 'Are you sure you want to delete this record?'
            this.modalFunction = function (data) {
              // let currentScope = this.$parent.$parent.$parent.$options.name
              let currentScope = this.$route.name
              axios.delete(config.BACKEND_API + currentScope.toLowerCase() + '/' + data.id).then(response => {
                this.$parent.$parent.$parent.$refs.vuetable.refresh()
              })
              .catch(e => {
                alert(e)
              })
            }
            this.record = data
            break
          default:
            console.log('nhan gi vay LAI')
        }
      }
    }
  }
</script>

<style>
  .custom-actions button.ui.button {
    padding: 8px 8px;
    box-shadow: none
  }
  .custom-actions button.ui.button > i.icon {
    margin: auto !important;
  }
</style>
