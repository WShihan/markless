function logout(){
    document.cookie = `marky-token=; expires=-1;path=/`;
    console.log('登出')
    window.location = '/login'
}