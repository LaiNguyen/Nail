<template>
  <div class="form-style">
    <div class="form-style-heading" v-if="expandMode == 'view'">Please View Detail Information</div>
    <div class="form-style-heading" v-else>Please Fill in Proper Detail Information</div>
    <label for="field1"><span>Name <span class="required">*</span></span>
      <input type="text" class="input-field" v-model="rowData.name" :value="rowData.name" :readonly="expandMode == 'view'"/>
    </label>
    <label for="field2"><span>Email <span class="required">*</span></span>
      <input type="text" class="input-field" v-model="rowData.email" :value="rowData.email" :readonly="expandMode == 'view'"/>
    </label>
    <label><span>Telephone</span>
      <input type="text" class="tel-number-field" value="" maxlength="3" v-model="tel_no_1" ref="phone_1"/>&nbsp;-
      <input type="text" class="tel-number-field" value="" maxlength="3" v-model="tel_no_2" ref="phone_2"/>&nbsp;-
      <input type="text" class="tel-number-field-long" value="" maxlength="4" v-model="tel_no_3" ref="phone_3"/>
    </label>
    <label for="field2"><span>Birthday</span>
      <datepicker class="input-field" v-model="rowData.birthday"></datepicker>
      </label>
    <label for="field4"><span>Gender</span>
      <div class="input-gender">&nbsp;
        <input type="radio" id="male" value="M" v-model="rowData.gender"><label for="male">&nbsp;Male&nbsp;&nbsp;&nbsp;&nbsp;</label>
        <input type="radio" id="female" value="F" v-model="rowData.gender"><label for="female">&nbsp;Female</label>
      </div>
    </label>
    <label for="field2"><span>Salary</span>
      <input type="number" class="input-field" v-model="rowData.salary" :value="rowData.salary" :readonly="expandMode == 'view'"/>
    </label>
    <label for="field2"><span>Commission</span>
      <input type="number" class="input-field" v-model="rowData.commission" :value="rowData.commission" :readonly="expandMode == 'view'"/>
    </label>
    <label for="field2"><span>Amount Own</span>
      <input type="number" class="input-field" v-model="rowData.amount_own" :value="rowData.amount_own" :readonly="expandMode == 'view'"/>
    </label>
    <!-- <label for="field5"><span>Message <span class="required">*</span></span><textarea name="field5" class="textarea-field"></textarea></label> -->
    <!-- <div class="form-style-heading"></div> -->
    <label v-if="expandMode == 'edit'"><span>&nbsp;</span>
      <input type="submit" value="Edit"  @click="onEdit"/>&nbsp;&nbsp;&nbsp;
      <!-- <input type="submit" value="Create"  @click="onCreate"/> -->
    </label>
  </div>
</template>

<script>
  import config from '@/config.js'
  import axios from 'axios'
  import Datepicker from 'vue-bulma-datepicker'
  export default {
    components: { Datepicker },
    data () {
      return {
        tel_no_1: null,
        tel_no_2: null,
        tel_no_3: null
      }
    },
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
        return this.$parent.$parent.expandMode
      }
    },
    methods: {
      onEdit (event) {
        let currentScope = this.$route.name
        axios.post(config.BACKEND_API + currentScope.toLowerCase() + '/' + this.rowData.id, {
          name: this.rowData.name,
          email: this.rowData.email,
          phone: this.rowData.phone,
          gender: this.rowData.gender,
          birthday: this.rowData.birthday,
          salary: parseFloat(this.rowData.salary),
          commission: parseFloat(this.rowData.commission),
          amount_own: parseFloat(this.rowData.amount_own),
          u_by: 'lai.nguyen@tripolis.com'
        }).then(response => {
          this.$parent.$parent.$refs.vuetable.toggleDetailRow(this.rowData.id)
          this.$parent.$parent.$refs.vuetable.refresh()
        })
        .catch(e => {
          alert(e)
        })
      },
      updatePhoneNumber () {
        if (this.tel_no_1 && this.tel_no_2 && this.tel_no_3) {
          this.rowData.phone = '(' + this.tel_no_1 + ')' + this.tel_no_2 + '-' + this.tel_no_3
        } else {
          this.rowData.phone = 'Please enter correct phone'
        }
      },
      devidePhoneNumber () {
        let firstPart = this.rowData.phone.split(')')
        this.tel_no_1 = firstPart[0].substring(1) // remove first character which is '('
        let secondPart = firstPart[1].split('-')
        this.tel_no_2 = secondPart[0]
        this.tel_no_3 = secondPart[1]
      }
    },
    beforeMount () {
      this.devidePhoneNumber()
    },
    watch: {
      tel_no_1 (newValue, oldValue) {
        if (newValue.length === 3) {
          this.$refs.phone_2.focus()
        }
        this.updatePhoneNumber()
      },
      tel_no_2 (newValue, oldValue) {
        if (newValue.length === 3) {
          this.$refs.phone_3.focus()
        } else if (newValue.length === 0) {
          this.$refs.phone_1.focus()
        }
        this.updatePhoneNumber()
      },
      tel_no_3 (newValue, oldValue) {
        if (newValue.length === 0) {
          this.$refs.phone_2.focus()
        }
        this.updatePhoneNumber()
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
      width: 100px;
      font-weight: bold;
      float: left;
      padding-top: 8px;
      padding-right: 5px;
  }
  .form-style .input-gender{
      display: flex;
      position: relative;
      top: 8px;
  }
  .form-style span.required{
      color:red;
  }
  .form-style .tel-number-field{
      width: 40px;
      text-align: center;
  }
  .form-style .tel-number-field-long{
      width: 60px;
      text-align: center;
  }
  .form-style input.input-field{
      width: 48%;
  }
  .form-style input.input-field, 
  .form-style .tel-number-field,
  .form-style .tel-number-field-long, 
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
  .form-style .tel-number-field-long:focus,
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