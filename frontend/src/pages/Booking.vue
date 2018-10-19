<template>
  <div class="hello">
    <full-calendar :events="fcEvents" locale="en"></full-calendar>
  </div>
</template>

<script>
import fullCalendar from 'vue-fullcalendar'
var demoEvents = [
  {
    title: 'Sunny Out of Office',
    start: '2017-12-02',
    end: '2017-12-02'
  }
]
export default {
  name: 'Booking',
  components: { fullCalendar },
  data () {
    return {
      fcEvents: demoEvents
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

a {
  color: #42b983;
}
</style>
