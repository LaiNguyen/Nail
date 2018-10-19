<template>
  <div class="custom-actions">
    <modal v-if="needToConfirm" :modalTitle="modalTitle" :modalFunction="modalFunction" :record="record"></modal>
    <checkout v-if="needToConfirmCheckout" :modalTitle="checkoutTitle" :modalFunction="modalFunction" :record="record"></checkout>
    <button class="ui basic button" @click="itemAction('edit', rowData, rowIndex)">
      <span class="ui orange label"><i class="large edit icon" style="margin: 0px"></i></span>
    </button>
    <button class="ui basic button" @click="itemAction('delete', rowData, rowIndex)">
      <span class="ui red label"><i class="large remove icon" style="margin: 0px"></i></span>
    </button>
    <button class="ui basic button" @click="itemAction('checkout', rowData, rowIndex)">
      <span class="ui green label"><i class="large money icon"></i>CHECK OUT</span>
    </button>
  </div>
</template>
<script>
  import config from '@/config.js'
  import axios from 'axios'
  import modal from '@/mixins/modal.vue'
  import checkout from '@/mixins/checkout.vue'
  export default {
    data () {
      return {
        needToConfirm: false,
        needToConfirmCheckout: false,
        modalTitle: '',
        checkoutTitle: '',
        record: {},
        modalFunction: function () {

        }
      }
    },
    components: { config, axios, modal, checkout },
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
          case 'checkout':
            this.needToConfirmCheckout = true
            this.checkoutTitle = data.customer.name + ' - checkout summary:'
            this.modalFunction = function (data, payment) {
              let currentScope = 'order'
              let staffTips = []
              for (var k in payment.staff_tips) {
                if (payment.staff_tips.hasOwnProperty(k)) {
                  staffTips.push({
                    staff_id: k,
                    tip_amount: parseFloat(payment.staff_tips[k].tip_amount)
                  })
                }
              }
              payment.staff_tips = staffTips
              axios.post(config.BACKEND_API + currentScope + '/' + data.id + '/checkout', {
                order: data,
                payment: payment
              }).then(response => {
                this.$parent.$parent.$parent.$refs.vuetable.refresh()
              })
              .catch(e => {
                alert(e)
              })
            }
            this.record = data
            break
          case 'edit':
            this.$parent.$parent.onEditClicked(data)
            break
          case 'delete':
            this.needToConfirm = true
            this.modalTitle = 'Are you sure you want to delete this record?'
            this.modalFunction = function (data) {
              let currentScope = 'order'
              axios.delete(config.BACKEND_API + currentScope + '/' + data.id).then(response => {
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
