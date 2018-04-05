function login(){
	user_id = document.getElementById('user_id')
	console.log(user_id)
	localStorage.setItem("user_id",user_id.value)
	window.location.href = "menu.html"
}