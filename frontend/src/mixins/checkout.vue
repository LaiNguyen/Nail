<template lang="html">
  <div class="modal">
    <div class="modal-window" style="overflow-y: scroll; max-height:750px;">
      <div>
        <vue-good-table
          :title="modalTitle"
          :columns="columns"
          :rows="calculateRows"
          :lineNumbers="true"/>
      </div>
      <div class="customer-payment">
        <h5>Customer Pay:<h1>${{total}}</h1></h5>
        <label><span>Surcharge:</span><input type="input" v-model="surcharge"/></label>
        <label><span>Discount:</span><input type="input" v-model="discount"/></label>
        <label><span>Payment:</span>
          <div class="input-payment">
            <input type="radio" id="cash" value="0" v-model="paymentType"/><label for="cash">&nbsp;Cash&nbsp;&nbsp;&nbsp;&nbsp;</label>
            <input type="radio" id="debit" value="1" v-model="paymentType"/><label for="debit">&nbsp;Debit&nbsp;&nbsp;&nbsp;&nbsp;</label>
            <input type="radio" id="credit" value="2" v-model="paymentType"/><label for="credit">&nbsp;Credit</label>
          </div>
        </label>
        <label><span>Tip type:</span>
          <div class="input-payment">
            <input type="radio" id="share" value="0" v-model="tipType"/><label for="share">&nbsp;Share&nbsp;&nbsp;&nbsp;&nbsp;</label>
            <input type="radio" id="specific" value="1" v-model="tipType"/><label for="specific">&nbsp;Specific</label>
          </div>
        </label>
        <label v-if="tipType==0"><span>Tip amount:</span><input type="input" v-model="tip"/></label>
        <div v-if="tipType==1" v-for="(item, index) in calculateStaffs">
        <label><span>Tip for {{item.name}}:</span><input type="input" v-model="item.tip_amount"/></label>
        <!-- <label><span>Tip for Lai:</span><input type="input" v-model="tip"/></label> -->
        </div>
      </div>
      <div class="actions">
       <button class="cancel" @click="cancel">Cancel</button>
       <button class="confirm" @click="confirm">Check Out</button>
      </div>
    </div>
  </div>
</template>
<script>
import VueGoodTable from 'vue-good-table'
import Vue from 'vue'
Vue.use(VueGoodTable)
export default {
  name: 'modal',
  props: ['modalTitle', 'modalFunction', 'record'],
  data () {
    return {
      columns: [
        {
          label: 'Service Used',
          field: 'name'
        },
        {
          label: 'Staff Name',
          field: 'staff_name'
        },
        {
          label: 'Price',
          field: 'price',
          type: 'number'
        }
      ],
      staffs: {},
      surcharge: 0,
      discount: 0,
      paymentType: 0,
      tipType: 0,
      tip: 0,
      total: this.record.total_price
    }
  },
  computed: {
    calculateRows () {
      var rows = []
      this.record.packages.forEach(function (item) {
        rows.push(
          {
            name: item.product_name,
            staff_name: item.staff_name,
            price: item.price
          }
        )
      })
      return rows
    },
    calculateStaffs () {
      var staffs = {}
      this.record.packages.forEach(function (item) {
        staffs[item.staff_id] = {
          name: item.staff_name,
          tip_amount: 0
        }
      })
      this.staffs = staffs
      return this.staffs
    }
  },
  methods: {
    modalShow (rec) {
      this.record = rec
    },
    confirm () {
      let payment = {
        surcharge: parseFloat(this.surcharge),
        discount: parseFloat(this.discount),
        payment_type: parseInt(this.paymentType),
        tip_type: parseInt(this.tipType),
        tip: parseFloat(this.tip),
        staff_tips: this.staffs
      }
      // this.record.total_price = this.total
      this.modalFunction(this.record, payment)
      this.$parent.needToConfirmCheckout = false
    },
    cancel () {
      this.$parent.needToConfirmCheckout = false
    },
    computeTotalWithSpecificTips () {
      var total = this.record.total_price
      if (this.surcharge) {
        total += parseInt(this.surcharge)
      }
      if (this.discount) {
        total -= parseInt(this.discount)
      }
      for (var k in this.staffs) {
        if (this.staffs.hasOwnProperty(k)) {
          if (this.staffs[k].tip_amount) {
            total += parseInt(this.staffs[k].tip_amount)
          }
        }
      }
      this.total = total
    }
  },
  watch: {
    tip (newValue, oldValue) {
      if (newValue) {
        this.total += parseInt(newValue)
      }
      if (oldValue) {
        this.total -= parseInt(oldValue)
      }
    },
    staffs: {
      handler: function () {
        this.computeTotalWithSpecificTips()
      },
      deep: true
    },
    tipType () {
      if (this.tipType === '0') {
        this.total = this.record.total_price
        if (this.tip) {
          this.total += parseInt(this.tip)
        }
        if (this.surcharge) {
          this.total += parseInt(this.surcharge)
        }
        if (this.discount) {
          this.total -= parseInt(this.discount)
        }
      } else {
        this.computeTotalWithSpecificTips()
      }
    },
    surcharge (newValue, oldValue) {
      if (newValue) {
        this.total += parseInt(newValue)
      }
      if (oldValue) {
        this.total -= parseInt(oldValue)
      }
    },
    discount (newValue, oldValue) {
      if (newValue) {
        this.total -= parseInt(newValue)
      }
      if (oldValue) {
        this.total += parseInt(oldValue)
      }
    }
  }
}
</script>
<style scoped>
  button {
   cursor: pointer;
  }

  .total {
    text-align: -webkit-auto;
    float: right;
    margin-right: 15px;
  }
  .customer-payment {
    text-align: -webkit-auto;
    margin-left: 11px;
    margin-top: 7px;
  }
  .customer-payment h5 {
    font-size: 18px;
  }
  .customer-payment input {
    box-sizing: border-box;
    -webkit-box-sizing: border-box;
    -moz-box-sizing: border-box;
    border: 1px solid #C2C2C2;
    box-shadow: 1px 1px 4px #EBEBEB;
    -moz-box-shadow: 1px 1px 4px #EBEBEB;
    -webkit-box-shadow: 1px 1px 4px #EBEBEB;
    border-radius: 3px;
    -webkit-border-radius: 3px;
    -moz-border-radius: 3px;
    padding: 7px;
    outline: none;
  }
  .customer-payment label{
      display: block;
      margin: 0px 0px 10px 0px;
  }
  .customer-payment label > span{
      min-width: 130px;
      width: auto;
      font-weight: bold;
      float: left;
      padding-top: 8px;
      padding-right: 5px;
  }
  .customer-payment .input-payment{
      display: -webkit-box;
      position: relative;
      top: 8px;
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
   top: 50%;
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
   /*text-align: center;*/
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
