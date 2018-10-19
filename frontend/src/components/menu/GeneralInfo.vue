<template>
  <div class="wrapper">
<!--     <div class="info">
      <div>Booking: 12</div>
      <div>Serving: 10</div>
    </div> -->
    <div class="logout">
      <span class="ui olive label" @click="logout"><i class="large sign out icon"></i>Sign Out</span>
    </div>
    <div class="linkedin">
      <a href="https://www.linkedin.com/in/lai-nguyen-260a88133/">
        <span class="ui label"><i class="large linkedin square icon"></i>Contact Us</span>
      </a>
    </div>
  </div>
</template>

<script>
import icon from 'vue-icon'
import config from '@/config.js'
import axios from 'axios'
export default {
  name: 'general-info',
  components: { icon },
  data () {
    return {
    }
  },
  methods: {
    logout () {
      axios.post(config.BACKEND_API + 'session/logout', {
      }).then(response => {
        localStorage.removeItem('access_token')
        location.reload()
      })
      .catch(e => {
        alert(e)
        localStorage.removeItem('access_token')
        location.reload()
      })
    }
  },
  watch: {
  }
}
</script>

<style scoped>
div .wrapper {
/*  width: 150px;
  height: 63px;
  position: relative;
  background-color: #4f8671;
  color: white;*/
}

div .info {
  margin-bottom: -1px;
  margin-left: 10px;
  line-height: 26px;
  font-size: 15px;
  font-weight: bold;
}

div .logout span{
  cursor: pointer;
  width: 100%;
}

.label {
  border-radius: 0 !important;
}

div .facebook span {
  background-color: #4468B0;
  border-color: #4468B0;
  color: #FFF;
  cursor: pointer;
  width: 100%
}
div .linkedin span {
  background-color: #0077B5;
  border-color: #0077B5;
  color: #FFF;
  cursor: pointer;
  width: 100%;
  height: 31px;
}
</style>
