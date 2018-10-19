<template>
  <div id="app">
    <Login v-if="!authenticated"></Login>
    <div v-else class="app">
      <div class="big-menu">
        <bigMenu :navigation="'home'" :label="'HOME'" :iconName="'home'" :currentModule="currentModule"/>
        <bigMenu :navigation="'booking'" :label="'BOOKING'" :iconName="'remind'" :currentModule="currentModule"/>
        <bigMenu :navigation="'service'" :label="'SERVICE'" :iconName="'apps'" :currentModule="currentModule"/>
        <bigMenu :navigation="'customer'" :label="'CUSTOMER'" :iconName="'attention'" :currentModule="currentModule"/>
        <bigMenu :navigation="'staff'" :label="'STAFF'" :iconName="'people'" :currentModule="currentModule"/>
        <bigMenu :navigation="'billing'" :label="'BILLING'" :iconName="'share'" :currentModule="currentModule"/>
        <bigMenu :navigation="'report'" :label="'REPORT'" :iconName="'form'" :currentModule="currentModule"/>
        <bigMenu :navigation="'giftcard'" :label="'GIFT CARD'" :iconName="'like'" :currentModule="currentModule"/>
        <bigMenu :navigation="'setting'" :label="'SETTING'" :iconName="'settings'" :currentModule="currentModule"/>
        <bigMenu :navigation="'review'" :label="'REVIEW'" :iconName="'appreciate'" :currentModule="currentModule"/>
        <generalInfo/>
      </div>
      <div class="content">
        <router-view/>  
      </div>
    </div>
  </div>
</template>

<script>
import bigMenu from '@/components/menu/BigMenu'
import generalInfo from '@/components/menu/GeneralInfo'
import Login from '@/components/auth/Login'
import auth from '@/components/auth'
import axios from 'axios'

axios.interceptors.request.use(config => {
// config.headers.Authorization = 'Bearer ' + window.localStorage.getItem('token')
  config.headers.Authorization = window.localStorage.getItem('access_token')
  config.headers.Accept = 'application/json'
  return config
})

export default {
  name: 'app',
  components: { bigMenu, generalInfo, Login, auth },
  data () {
    return {
      currentModule: null
    }
  },
  methods: {
    updateCurrentModule (value) {
      this.currentModule = value
    }
  },
  computed: {
    authenticated () {
      return auth.checkAuth()
    }
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  /*text-align: center;*/
  color: #2c3e50;
  /*margin-top: 60px;*/
  /*display: flex;*/
  /*margin: -8px;*/
  font-weight: 500;
  letter-spacing: 1px;
}

h1 {
  font-family: 'Vibur', cursive;
}

.app {
  display: flex;
}
.content {
  width: 100%;
  /*margin-right: 10px;*/
  /*margin-left: 10px;*/
  /*margin: 10px;*/
}
.big-menu {
  display: block;
}
.big-actions {
  width: 100%;
  position: absolute;
  left: 150px;
  top: 7px;
}
.big-actions .label {
  letter-spacing: 5px;
  font-size: 15px;
}
.big-actions button.ui.button {
  padding: 8px 8px;
  box-shadow: none !important;
}
.big-actions button.ui.button:focus{
  outline: 0;
}
</style>
