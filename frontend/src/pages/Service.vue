<template>
  <div class="service">
    <!-- <h1>Service Page</h1> -->
    <div class="big-actions">
      <button class="ui basic button" @click="onAdd">
        <span class="ui teal label"><i class="large add square icon"></i>ADD SERVICE</span>
      </button>
    </div>
    <template>
      <global-table
        :api-url="apiURL"
        :fields="fields"
        :sort-order="sortOrder"
        :append-params="moreParams"
        detail-row-component="detail-row-input-service"
      ></global-table>
    </template>
  </div>
</template>
        <!-- api-url="https://vuetable.ratiw.net/api/users" -->
<script>

import GlobalTable from '@/components/table/GlobalTable'
import ServiceDetailInput from '@/components/table/ServiceDetailInput'
import config from '@/config.js'
import axios from 'axios'
import Vue from 'vue'
Vue.component('detail-row-input-service', ServiceDetailInput)

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
          name: 'name',
          sortField: 'name'
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
          field: 'c_at',
          sortField: 'c_at',
          direction: 'desc'
        }
      ],
      moreParams: {},
      apiURL: config.BACKEND_API + 'service/?paging=yes'
    }
  },
  methods: {
    onCellClicked (data) {
      this.$router.push({path: 'product', query: { service_id: data.id, service_name: data.name }})
    },
    onAdd () {
      this.sortOrder[0].field = 'c_at'
      this.sortOrder[0].sortField = 'c_at'
      this.sortOrder[0].direction = 'desc'
      axios.post(config.BACKEND_API + 'service/', {
        service_id: this.$route.query.service_id,
        name: 'New Service',
        price: 0,
        duration: 0,
        c_by: 'lai.nguyen@tripolis.com'
      }).then(response => {
        this.$children[0].$refs.vuetable.refresh()
      })
      .catch(e => {
        alert(e)
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  h1 {
    color: #6daa94;
    text-align: center;
    /* font-family: 'Vibur', fantasy; */
    font-size: 30px;
    font-weight: bold;
    margin-bottom: -55px;
    width: 220px;
    border-top: 2px solid #8f8f8f;
    border-bottom: 2px solid #8f8f8f;
  }
</style>
