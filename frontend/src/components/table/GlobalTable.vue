<template>
  <div class="global-table">
    <template>
      <div class="loading">
        <sync-loader :loading="loading" :color="color" :size="size"></sync-loader>
      </div>
      <div v-show="!loading">
        <div v-show="displayFilter">
          <filter-bar></filter-bar>
        </div>
        <vuetable ref="vuetable"
          :api-url="apiUrl"
          :fields="fields"
          :detail-row-component="getDetailRow"
          detail-row-transition="fade"
          pagination-path=""
          :per-page="8"
          :sort-order="sortOrder"
          :append-params="appendParams"
          @vuetable:cell-clicked="onCellClicked"
          @vuetable:pagination-data="onPaginationData"
          @vuetable:loaded="onDataLoaded">
        </vuetable>
        <div class="vuetable-pagination ui basic segment grid">
          <vuetable-pagination-info ref="paginationInfo"></vuetable-pagination-info>
          <vuetable-pagination ref="pagination" @vuetable-pagination:change-page="onChangePage"></vuetable-pagination>
        </div>
      </div>
    </template>
  </div>
</template>

<script>

import Vuetable from 'vuetable-2/src/components/Vuetable'
import VuetablePagination from 'vuetable-2/src/components/VuetablePagination'
import VuetablePaginationInfo from 'vuetable-2/src/components/VuetablePaginationInfo'
// import VuetablePagination from 'vuetable-2/src/components/VuetablePaginationDropdown'

import VueEvents from 'vue-events'
import SyncLoader from 'vue-spinner/src/SyncLoader.vue'

import CustomActions from '@/components/action/CustomActions'
import HomeActions from '@/components/action/HomeActions'
import FilterBar from '@/components/table/FilterBar'

import accounting from 'accounting'
import moment from 'moment'
import Vue from 'vue'

Vue.component('custom-actions', CustomActions)
Vue.component('home-actions', HomeActions)
Vue.use(VueEvents)

export default {
  name: 'global-table',
  components: { Vuetable, VuetablePagination, VuetablePaginationInfo, FilterBar, VueEvents, SyncLoader },
  data () {
    return {
      loading: true,
      expandMode: 'view'
    }
  },
  props: {
    apiUrl: {
      type: String,
      required: true
    },
    fields: {
      type: Array,
      required: true
    },
    sortOrder: {
      type: Array,
      default () {
        return []
      }
    },
    appendParams: {
      type: Object,
      default () {
        return {}
      }
    },
    color: '#3AB982',
    size: {
      height: '35px',
      width: '4px',
      margin: '2px',
      radius: '2px'
    },
    displayFilter: {
      type: Boolean,
      default () {
        return true
      }
    },
    detailRowComponent: {
      type: String,
      default () {
        return {}
      }
    }
  },
  computed: {
    getDetailRow () {
      return this.detailRowComponent
    }
  },
  mounted () {
    this.$events.$on('filter-set', eventData => this.onFilterSet(eventData))
    this.$events.$on('filter-reset', e => this.onFilterReset())
  },
  methods: {
    allcap (value) {
      return value.toUpperCase()
    },
    genderLabel (value) {
      return value === 'M'
        ? '<span class="ui teal label"><i class="large man icon"></i>Male</span>'
        : '<span class="ui pink label"><i class="large woman icon"></i>Female</span>'
        // : '<span class="ui pink label"><i class="large female icon" style="margin: 0px"></i></span>'
    },
    favoriteLabel (value) {
      if (value === 0) {
        return '<span class="ui yellow label"><i class="large child icon"></i>New</span>'
      } else if (value === 1) {
        return '<span class="ui yellow label"><i class="large student icon"></i>Normal</span>'
      } else if (value === 2) {
        return '<span class="ui yellow label"><i class="large history icon"></i>Regular</span>'
      }
      return '<span class="ui yellow label"><i class="large star icon"></i>VIP</span>'
        // : '<span class="ui pink label"><i class="large female icon" style="margin: 0px"></i></span>'
    },
    paymentLabel (value) {
      if (value === 0) {
        return 'CASH'
      } else if (value === 1) {
        return 'DEBIT'
      }
      return 'CREDIT'
    },
    billingStatusLabel (value) {
      if (value === 0) {
        return 'Unpaid'
      } else if (value === 1) {
        return 'Paid'
      }
      return 'Refund'
    },
    orderNumberLabel (value) {
      return value.substr(value.length - 6)
    },
    formatNumber (value) {
      return '$' + accounting.formatNumber(value, 2)
    },
    formatDuration (value) {
      return value + ' minutes'
    },
    formatBirthday (value, fmt = 'D MMM YYYY') {
      return (value == null)
        ? ''
        : moment(value, 'YYYY-MM-DD').format(fmt)
    },
    formatDate (value, fmt = 'D MMM YYYY') {
      return (value == null)
        ? ''
        : moment(value).format(fmt)
    },
    onPaginationData (paginationData) {
      this.$refs.pagination.setPaginationData(paginationData)
      this.$refs.paginationInfo.setPaginationData(paginationData)
    },
    onChangePage (page) {
      this.$refs.vuetable.changePage(page)
    },
    onCellClicked (data, field, event) {
      console.log('onCellClicked')
      this.$parent.onCellClicked(data)
    },
    onViewClicked (data) {
      this.expandMode = 'view'
      // this.$refs.vuetable.hideDetailRow(data.id)
      this.$refs.vuetable.toggleDetailRow(data.id)
    },
    onEditClicked (data) {
      this.expandMode = 'edit'
      // this.$refs.vuetable.hideDetailRow(data.id)
      this.$refs.vuetable.toggleDetailRow(data.id)
    },
    onFilterSet (filterText) {
      this.appendParams.filter = filterText
      Vue.nextTick(() => this.$refs.vuetable.refresh())
    },
    onFilterReset () {
      delete this.appendParams.filter
      Vue.nextTick(() => this.$refs.vuetable.refresh())
    },
    onDataLoaded () {
      this.loading = false
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.loading {
  /*display: flex;
  justify-content: center;*/
  position: absolute;
  top: 300px;
  right: 670px;
}
.filter-bar {
  float: right;
}
</style>
