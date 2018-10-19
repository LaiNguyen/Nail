<template>
  <div class="customer">
    <!-- <h1>Billing Page</h1> -->
    <template>
      <global-table
        :api-url="apiURL"
        :fields="fields"
        :sort-order="sortOrder"
        :append-params="moreParams"
        detail-row-component="detail-row-input-billing"
      ></global-table>
    </template>
  </div>
</template>
        <!-- api-url="https://vuetable.ratiw.net/api/users" -->
<script>

import GlobalTable from '@/components/table/GlobalTable'
import BillingDetailInput from '@/components/table/BillingDetailInput'
import config from '@/config.js'
// import axios from 'axios'

import Vue from 'vue'
Vue.component('detail-row-input-billing', BillingDetailInput)

export default {
  name: 'Customer',
  components: { GlobalTable },
  data () {
    return {
      fields: [
        {
          name: '__sequence',
          title: '#',
          titleClass: 'center aligned',
          dataClass: 'center aligned'
        },
        {
          title: 'Order No.',
          name: 'order_id',
          callback: 'orderNumberLabel'
        },
        {
          title: 'Initial',
          name: 'initial_amount',
          sortField: 'initial_amount'
        },
        {
          name: 'surcharge',
          sortField: 'surcharge'
        },
        {
          name: 'discount',
          sortField: 'discount'
        },
        {
          name: 'tip',
          sortField: 'tip'
        },
        {
          title: 'Total',
          name: 'total_amount',
          sortField: 'total_amount'
        },
        {
          name: 'revenue',
          sortField: 'revenue'
        },
        {
          title: 'Payment',
          name: 'payment_type',
          sortField: 'payment_type',
          callback: 'paymentLabel'
        },
        {
          name: 'status',
          sortField: 'status',
          callback: 'billingStatusLabel'
        },
        {
          title: 'Created At',
          name: 'c_at',
          sortField: 'c_at',
          titleClass: 'center aligned',
          dataClass: 'center aligned',
          callback: 'formatDate|MMMM Do YYYY, hh:mm a'
        },
        {
          name: '__component:custom-actions',
          title: 'Actions',
          titleClass: 'center aligned',
          dataClass: 'center aligned'
        }
      ],
      sortOrder: [
        {
          field: 'name',
          sortField: 'name',
          direction: 'asc'
        }
      ],
      moreParams: {},
      apiURL: config.BACKEND_API + 'billing/?paging=yes'
    }
  },
  methods: {
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  h1 {
    color: #6daa94;
    text-align: center;
    font-size: 30px;
    font-weight: bold;
    margin-bottom: -55px;
    width: 220px;
    border-top: 2px solid #8f8f8f;
    border-bottom: 2px solid #8f8f8f;
  }
</style>
