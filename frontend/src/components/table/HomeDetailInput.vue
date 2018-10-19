<template>
  <div class="form-style">
    <div class="form-style-heading">Customer information:</div>
    <label for="field2"><span>Customer</span>
      <multiselect
        v-model="selectedCustomer"
        :options="allCustomers"
        :multiple="false"
        track-by="name"
        :custom-label="customLabel"
        placeholder="Select Customer">
      </multiselect>
    </label>
    <div class="form-style-heading">Customer {{customerName}} wants to use:</div>
    <div v-for="(item, index) in packages">
      <label for="field1"><span>Product</span>
        <ProductRow
          :index="index"
          :package="item"
          :staffs="allStaffs"
          :services="allServices"
          :selectedStaffID="item.staff_id"
          :selectedServiceID="item.service_id"
          :selectedProductID="item.product_id"></ProductRow>
      </label>
    </div>
    <div class="actions">
      <button class="ui basic button" @click="onAdd">
        <span class="ui green label"><i class="large add square icon"></i>ADD PRODUCT</span>
      </button>
      <button class="ui basic button" @click="onSave">
        <span class="ui green label"><i class="large save icon"></i>SAVE ORDERS</span>
      </button>
    </div>
  </div>
</template>

<script>
  import config from '@/config.js'
  import axios from 'axios'
  import Multiselect from 'vue-multiselect'
  import ProductRow from '@/components/table/ProductRow'
  import _ from 'lodash'
  export default {
    data () {
      return {
        customerName: this.rowData.customer.name,
        selectedCustomer: null
      }
    },
    components: { config, axios, Multiselect, ProductRow },
    props: {
      rowData: {
        type: Object,
        required: true
      },
      rowIndex: {
        type: Number
      }
    },
    computed: {
      phoneFormat (value) {
        return 'ABC'
        // return value.replace(/[^0-9]/g, '').replace(/(\d{3})(\d{3})(\d{4})/, '($1) $2-$3')
      },
      allCustomers () {
        return this.$parent.$parent.$parent.customers
      },
      allStaffs () {
        return this.$parent.$parent.$parent.staffs
      },
      allServices () {
        return this.$parent.$parent.$parent.services
      },
      packages () {
        return this.rowData.packages
      }
    },
    methods: {
      onSave (event) {
        axios.post(config.BACKEND_API + 'order/' + this.rowData.id, {
          customer_id: this.rowData.customer.id,
          packages: this.rowData.packages,
          total_price: this.rowData.total_price,
          total_duration: this.rowData.total_duration,
          status: this.rowData.status,
          started_at: this.rowData.started_at,
          ended_at: this.rowData.ended_at,
          u_by: 'lai.nguyen@tripolis.com'
        }).then(response => {
          this.$parent.$parent.$refs.vuetable.toggleDetailRow(this.rowData.id)
          this.$parent.$parent.$refs.vuetable.refresh()
        })
        .catch(e => {
          alert(e)
        })
      },
      onAdd (event) {
        let newPackage = {
        }
        this.rowData.packages.push(newPackage)
      },
      customLabel (option) {
        return `${option.name}`
      },
      genderLabel (value) {
        return value === 'M' ? 'Male' : 'Female'
      },
      removePackage (index) {
        this.rowData.packages.splice(index, 1)
        this.updatePackages()
      },
      updatePackages () {
        var newRowData = this.rowData
        newRowData.total_price = 0
        newRowData.total_duration = 0
        newRowData.packages.forEach(function (item) {
          newRowData.total_price += item.price
          newRowData.total_duration += item.duration
        })
        this.rowData = newRowData
      },
      calculateSelectedCustomer () {
        this.selectedCustomer = _.find(this.allCustomers, {id: this.rowData.customer.id})
      }
    },
    watch: {
      selectedCustomer: function () {
        this.rowData.customer = this.selectedCustomer
        this.customerName = this.rowData.customer.name
      }
    },
    beforeMount () {
      if (this.rowData.customer.id) {
        this.calculateSelectedCustomer()
      }
    }
  }
</script>
<style src="vue-multiselect/dist/vue-multiselect.min.css"></style>
<style scoped>
  .actions button.ui.button {
    padding: 8px 8px;
    box-shadow: none;
  }
  .actions {
    margin-top: -15px;
  }
  .actions button.ui.button > i.icon {
    margin: auto !important;
  }
  .multiselect{
    left: 72px;
    top: -28px;
    width: 300px;
  }
  .form-style{
      /*max-width: 1100px;*/
      padding: 20px 12px 10px 20px;
      font: 13px Arial, Helvetica, sans-serif;
  }
  .form-style-heading{
      font-weight: bold;
      font-style: italic;
      border-bottom: 2px solid #ddd;
      margin-bottom: 30px;
      font-size: 15px;
      padding-bottom: 10px;
  }
  .form-style label{
      display: block;
  }
  .form-style label > span{
      width: 65px;
      font-weight: bold;
      /*float: left;*/
      padding-top: 8px;
      padding-right: 5px;
  }
  .form-style .label-select{
      padding-top: 2px;
  }
  .form-style span.required{
      color:red;
  }
  .form-style .tel-number-field{
      width: 40px;
      text-align: center;
  }
  .form-style input.input-field{
      width: 48%; 
  }

  .form-style input.input-field, 
  .form-style .tel-number-field, 
  .form-style .textarea-field, 
  .form-style .select-field{
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
  .form-style .input-field:focus, 
  .form-style .tel-number-field:focus, 
  .form-style .textarea-field:focus,  
  .form-style .select-field:focus{
      border: 1px solid #0C0;
  }
  .form-style .textarea-field{
      height:100px;
      width: 55%;
  }
  .form-style .checkbox-field{
      cursor: pointer;
      transform:scale(1.4, 1.4);
      position: relative;
      top: 6px;
      left: 3px;
  }
  .form-style input[type=submit],
  .form-style input[type=button]{
      width: 70px;
      border: none;
      padding: 8px 15px 8px 15px;
      background: #FF8500;
      color: #fff;
      box-shadow: 1px 1px 4px #DADADA;
      -moz-box-shadow: 1px 1px 4px #DADADA;
      -webkit-box-shadow: 1px 1px 4px #DADADA;
      border-radius: 3px;
      -webkit-border-radius: 3px;
      -moz-border-radius: 3px;
  }
  .form-style input[type=submit]:hover,
  .form-style input[type=button]:hover{
      background: #EA7B00;
      color: #fff;
  }
  </style>