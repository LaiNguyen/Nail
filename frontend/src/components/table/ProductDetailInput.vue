<template>
  <div class="form-style">
    <div class="form-style-heading" v-if="expandMode == 'view'">Please View Detail Information</div>
    <div class="form-style-heading" v-else>Please Fill in Product Detail Information</div>
    <label for="field1"><span>Service Name </span>
      <input type="text" class="input-field" v-model="serviceName" :value="serviceName" readonly="true" style="background: transparent;border: none"/>
    </label>
    <label for="field1"><span>Product Name <span class="required">*</span></span>
      <input type="text" class="input-field" v-model="rowData.name" :value="rowData.name" :readonly="expandMode == 'view'"/>
    </label>
    <label for="field2"><span>Price</span>
      <input type="text" class="input-field" v-model="rowData.price" :value="rowData.price" :readonly="expandMode == 'view'"/>
    </label>
    <label for="field2"><span>Duration</span>
      <input type="text" class="input-field" v-model="rowData.duration" :value="rowData.duration" :readonly="expandMode == 'view'"/>
    </label>

    <label v-if="expandMode == 'edit'"><span>&nbsp;</span>
      <input type="submit" value="Edit"  @click="onEdit"/>
    </label>
  </div>
</template>

<script>
  import config from '@/config.js'
  import axios from 'axios'
  export default {
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
      expandMode () {
        return this.$parent.$parent.expandMode // global table
      },
      serviceName () {
        return this.$parent.$parent.$parent.serviceName // product vue
      }
    },
    methods: {
      onEdit (event) {
        let currentScope = this.$route.name
        axios.post(config.BACKEND_API + currentScope.toLowerCase() + '/' + this.rowData.id, {
          service_id: this.$route.query.service_id,
          name: this.rowData.name,
          price: parseInt(this.rowData.price),
          duration: parseInt(this.rowData.duration),
          u_by: 'lai.nguyen@tripolis.com'
        }).then(response => {
          this.$parent.$parent.$refs.vuetable.toggleDetailRow(this.rowData.id)
          this.$parent.$parent.$refs.vuetable.refresh()
        })
        .catch(e => {
          alert(e)
        })
      }
    }
  }
  </script>
  <style scoped>
  .form-style{
      max-width: 500px;
      padding: 20px 12px 10px 20px;
      font: 13px Arial, Helvetica, sans-serif;
  }
  .form-style-heading{
      font-weight: bold;
      font-style: italic;
      border-bottom: 2px solid #ddd;
      margin-bottom: 20px;
      font-size: 15px;
      padding-bottom: 3px;
  }
  .form-style label{
      display: block;
      margin: 0px 0px 15px 0px;
  }
  .form-style label > span{
      width: 122px;
      font-weight: bold;
      float: left;
      padding-top: 8px;
      padding-right: 5px;
  }
  .form-style .label-gender{
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