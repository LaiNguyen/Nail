<template>
  <div class="product-row">
    <div class="select-service">
      <multiselect
        v-model="selectedService"
        :options="allServices"
        :multiple="false"
        track-by="name"
        :custom-label="customLabel"
        placeholder="Select Service">
      </multiselect>
    </div>
    <div class="select-product">
      <multiselect
        v-model="selectedProduct"
        :options="allProducts"
        :multiple="false"
        track-by="name"
        :custom-label="customProductLabel"
        placeholder="Select Product">
      </multiselect>
    </div>
    <div class="select-staff">
      <multiselect
        v-model="selectedStaff"
        :options="allStaffs"
        :multiple="false"
        track-by="name"
        :custom-label="customLabel"
        placeholder="Select Staff">
      </multiselect>
    </div>
    <div class="delete-button">
      <button class="ui basic button" @click="remove()">
        <span class="ui red label"><i class="large remove icon" style="margin: 0px"></i></span>
      </button>
    </div>
  </div>
</template>

<script>
  import Multiselect from 'vue-multiselect'
  import config from '@/config.js'
  import axios from 'axios'
  import _ from 'lodash'
  export default {
    data () {
      return {
        selectedStaff: null,
        selectedService: null,
        selectedProduct: null,
        products: null
      }
    },
    components: { Multiselect },
    props: {
      index: null,
      package: null,
      staffs: null,
      services: null,
      selectedStaffID: null,
      selectedServiceID: null,
      selectedProductID: null
    },
    computed: {
      allStaffs () {
        return this.staffs
      },
      allServices () {
        return this.services
      },
      allProducts () {
        if (this.products == null) {
          return []
        }
        return this.products
      }
    },
    methods: {
      customLabel (option) {
        return `${option.name}`
      },
      customProductLabel (option) {
        return `${option.name} - $${option.price} - ${option.duration} minutes`
      },
      fetchProduct () {
        axios.get(config.BACKEND_API + 'product/find_by_service/' + this.selectedService.id).then(response => {
          this.products = response.data.data
          this.calculateSelectedProduct()
        })
        .catch(e => {
          alert(e)
        })
      },
      calculateSelectedStaff () {
        this.selectedStaff = _.find(this.staffs, {id: this.selectedStaffID})
      },
      calculateSelectedService () {
        this.selectedService = _.find(this.services, {id: this.selectedServiceID})
      },
      calculateSelectedProduct () {
        this.selectedProduct = _.find(this.products, {id: this.selectedProductID})
      },
      remove () {
        this.$parent.removePackage(this.index)
      }
    },
    beforeMount () {
      if (this.selectedStaffID) {
        this.calculateSelectedStaff()
      }
      if (this.selectedServiceID) {
        this.calculateSelectedService()
        this.fetchProduct()
      }
    },
    watch: {
      selectedStaff: function () {
        this.package.staff_id = this.selectedStaff.id
        this.package.staff_name = this.selectedStaff.name
      },
      selectedService: function () {
        this.fetchProduct()
        this.package.service_id = this.selectedService.id
        this.package.service_name = this.selectedService.name
      },
      selectedProduct: function () {
        if (this.selectedProduct != null) {
          this.package.product_id = this.selectedProduct.id
          this.package.product_name = this.selectedProduct.name
          this.package.price = this.selectedProduct.price
          this.package.duration = this.selectedProduct.duration
          // console.log(this.package)
          this.$parent.updatePackages()
        }
      }
    }
  }
</script>
<style src="vue-multiselect/dist/vue-multiselect.min.css"></style>
<style scoped>
  button.ui.button {
    padding: 8px 8px;
    box-shadow: none
  }
  button.ui.button > i.icon {
    margin: auto !important;
  }

  .delete-button{
    position: relative;
    height: 0px;
    top: -32px;
    left: 60px;
  }

  .product-row{
    display: flex;
  }
  .multiselect{
    left: 60px;
    top: -28px;
  }
  .select-staff{
    width: 300px;
  }
  .select-service{
    width: 300px;
  }
  .select-product{
    width: 450px;
  }
  </style>