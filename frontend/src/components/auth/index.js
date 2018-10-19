// URL and endpoint constants
import axios from 'axios'
// const API_URL = 'http://138.128.65.43:4011/'
const API_URL = 'http://localhost:4011/'
const LOGIN_URL = API_URL + 'session/login'
const SIGNUP_URL = API_URL + 'users/'

export default {

  // User object will let us check authentication status
  user: {
    authenticated: false
  },

  // Send a request to the login URL and save the returned JWT
  login (context, creds, redirect) {
    axios.post(LOGIN_URL, creds).then(response => {
      // localStorage.setItem('id_token', response.data.id_token)
      localStorage.setItem('access_token', response.data.token + '|' + response.data.tenant_id)

      this.user.authenticated = true

      // Redirect to a specified route
      location.reload()
    })
    .catch(e => {
      alert(e)
    })
  },

  signup (context, creds, redirect) {
    context.$http.post(SIGNUP_URL, creds, (data) => {
      // localStorage.setItem('id_token', data.id_token)
      localStorage.setItem('access_token', data.access_token)

      this.user.authenticated = true

      if (redirect) {
        // router.go(redirect)
        this.$router.push({ path: redirect })
      }
    }).error((err) => {
      context.error = err
    })
  },

  // To log out, we just need to remove the token
  logout () {
    // localStorage.removeItem('id_token')
    localStorage.removeItem('access_token')
    this.user.authenticated = false
  },

  checkAuth () {
    var jwt = localStorage.getItem('access_token')
    if (jwt) {
      this.user.authenticated = true
    } else {
      this.user.authenticated = false
    }
    return this.user.authenticated
  },

  // The object to be passed as a header for authenticated requests
  getAuthHeader () {
    return {
      'Authorization': 'Bearer ' + localStorage.getItem('access_token')
    }
  }
}
