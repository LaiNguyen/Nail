<template>
  <div class="home">
    <div class="top-menu">
      <!-- <h1>Home Page</h1> -->
      <div class="big-actions">
        <button class="ui basic button" @click="onAdd">
          <span class="ui teal label"><i class="large add square icon"></i>ADD WALK IN</span>
        </button>
      </div>
    </div>
    <template>
      <global-table
        :api-url="apiURL"
        :fields="fields"
        :sort-order="sortOrder"
        :append-params="moreParams"
        :display-filter="displayFilter"
        detailRowComponent="detail-row-input-home"
      ></global-table>
    </template>
    <!-- <div class="upper-area">
      <div style="overflow-y: scroll; height:400px;">
        asd
      </div>
      <div class="order-detail">
        ysz
      </div>
    </div>
    <div class="lower-area">
      <div class="service">
        iop
      </div>
      <div class="product">
        bnm
      </div>
    </div> -->
  </div>
</template>

<script>
import config from '@/config.js'
import axios from 'axios'
import GlobalTable from '@/components/table/GlobalTable'
import HomeDetailInput from '@/components/table/HomeDetailInput'
import Vue from 'vue'

Vue.component('detail-row-input-home', HomeDetailInput)
export default {
  name: 'HOME',
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
          name: 'customer.name',
          sortField: 'customer.name',
          title: 'Customer'
        },
        {
          name: 'customer.phone',
          sortField: 'customer.phone',
          title: 'Phone'
        },
        {
          name: 'customer.status',
          sortField: 'customer.status',
          title: 'Type',
          titleClass: 'center aligned',
          dataClass: 'center aligned',
          callback: 'favoriteLabel'
        },
        {
          name: 'started_at',
          sortField: 'started_at',
          title: 'Check In Time',
          titleClass: 'center aligned',
          dataClass: 'center aligned',
          callback: 'formatDate|MMMM Do YYYY, hh:mm a'
        },
        {
          name: 'total_price',
          title: 'Total Amount',
          callback: 'formatNumber'
        },
        {
          name: 'total_duration',
          title: 'Total Duration',
          titleClass: 'center aligned',
          dataClass: 'center aligned',
          callback: 'formatDuration'
        },
        {
          name: '__component:home-actions',
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
      apiURL: config.BACKEND_API + 'order/?status=1',
      customers: [],
      staffs: [],
      services: []
    }
  },
  methods: {
    allCustomers () {
      axios.get(config.BACKEND_API + 'customer/').then(response => {
        this.customers = response.data.data
      })
      .catch(e => {
        alert(e)
      })
    },
    allStaffs () {
      axios.get(config.BACKEND_API + 'staff/').then(response => {
        this.staffs = response.data.data
      })
      .catch(e => {
        alert(e)
      })
    },
    allServices () {
      axios.get(config.BACKEND_API + 'service/').then(response => {
        this.services = response.data.data
      })
      .catch(e => {
        alert(e)
      })
    },
    onAdd () {
      axios.post(config.BACKEND_API + 'order/', {
        total_price: 0,
        total_duration: 0,
        c_by: 'lai.nguyen@tripolis.com'
      }).then(response => {
        this.$children[0].$refs.vuetable.refresh()
      })
      .catch(e => {
        alert(e)
      })
    }
  },
  computed: {
    displayFilter () {
      return true
    }
  },
  beforeMount () {
    this.allCustomers()
    this.allStaffs()
    this.allServices()
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
.upper-area {
  display: flex;
  height: 50%;
}
.lower-area {
  display: flex;
  height: 50%;
}
.my-scrollbar{
  width: 35%;
  min-width: 300px;
  max-height: 450px;
}
.scroll-me{
  min-width: 750px;
}
</style>
