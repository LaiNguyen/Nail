<template>
  <div class="staff">
    <!-- <h1>Staff Page</h1> -->
    <div class="big-actions">
      <button class="ui basic button" @click="onAdd">
        <span class="ui teal label"><i class="large add square icon"></i>ADD STAFF</span>
      </button>
    </div>
    <template>
      <global-table
        :api-url="apiURL"
        :fields="fields"
        :sort-order="sortOrder"
        :append-params="moreParams"
        detail-row-component="detail-row-input-staff"
      ></global-table>
    </template>
  </div>
</template>

<script>

import GlobalTable from '@/components/table/GlobalTable'
import StaffDetailInput from '@/components/table/StaffDetailInput'
import config from '@/config.js'
import axios from 'axios'
import Vue from 'vue'
Vue.component('detail-row-input-staff', StaffDetailInput)

export default {
  name: 'Staff',
  components: { GlobalTable },
  data () {
    return {
      fields: [
        {
          name: '__sequence',   // <----
          title: '#',
          titleClass: 'center aligned',
          dataClass: 'center aligned'
        },
        // {
        //   name: '__handle',   // <----
        //   dataClass: 'center aligned'
        // },
        {
          name: 'name',
          sortField: 'name'
        },
        {
          name: 'email',
          sortField: 'email'
        },
        {
          name: 'phone'
        },
        {
          name: 'birthday',
          sortField: 'birthday',
          titleClass: 'center aligned',
          dataClass: 'center aligned',
          callback: 'formatDate|DD-MM-YYYY'
        },
        {
          name: 'gender',
          sortField: 'gender',
          titleClass: 'center aligned',
          dataClass: 'center aligned',
          callback: 'genderLabel',
          visible: true
        },
        {
          name: 'salary',
          sortField: 'salary',
          titleClass: 'center aligned',
          dataClass: 'center aligned',
          callback: 'formatNumber'
        },
        {
          name: 'commission',
          sortField: 'commission',
          titleClass: 'center aligned',
          dataClass: 'center aligned',
          callback: 'formatNumber'
        },
        {
          name: 'amount_own',
          sortField: 'amount_own',
          title: 'Amount Own',
          titleClass: 'center aligned',
          dataClass: 'center aligned',
          callback: 'formatNumber'
        },
        {
          name: '__component:custom-actions',   // <----
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
      apiURL: config.BACKEND_API + 'staff/?paging=yes'
    }
  },
  methods: {
    onAdd () {
      this.sortOrder[0].field = 'c_at'
      this.sortOrder[0].sortField = 'c_at'
      this.sortOrder[0].direction = 'desc'
      axios.post(config.BACKEND_API + 'staff/', {
        service_id: this.$route.query.service_id,
        name: 'New Staff',
        amount_own: 0,
        salary: 0,
        commission: 0,
        status: 0,
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
