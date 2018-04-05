var user_id = localStorage.getItem("user_id")

function writeName(){
	var uid = document.getElementById('uid')
	uid.innerHTML=user_id
}
window.onload = writeName

function add(){
	amount = document.getElementById('add_amount').value
	var r = new XMLHttpRequest();
	r.open("POST", "/users", true)
	r.setRequestHeader("Content-type", "application/json")
	r.send(JSON.stringify({TransactionNum: 1, UserId: user_id, Amount: parseInt(amount)}))
	r.onreadystatechange = function(){
		if(r.readyState == 4){
			alert(r.response)
		}
	}
}

function quote() {
	stock = document.getElementById('quote_stock').value
	var r = new XMLHttpRequest();
	r.open("POST", "/stocks/"+stock+"/quote", true)
	r.setRequestHeader("Content-type", "application/json")
	r.send(JSON.stringify({TransactionNum: 1, UserId: user_id, StockSymbol: stock}))
	r.onreadystatechange = function(){
		if(r.readyState == 4){
			alert(r.response)
		}
	}
}

function buy() {
	stock = document.getElementById('buy_stock').value
	amount = document.getElementById('buy_amount').value
	var r = new XMLHttpRequest();
	r.open("POST", "/stocks/"+stock+"/buy", true)
	r.setRequestHeader("Content-type", "application/json")
	r.send(JSON.stringify({TransactionNum: 1, UserId: user_id, StockSymbol: stock, Amount: parseInt(amount)}))
	r.onreadystatechange = function(){
		if(r.readyState == 4){
			alert(r.response)
		}
	}
}

function commitbuy() {
	var r = new XMLHttpRequest();
	r.open("PUT", "/stocks/buy", true)
	r.setRequestHeader("Content-type", "application/json")
	r.send(JSON.stringify({TransactionNum: 1, UserId: user_id}))
	r.onreadystatechange = function(){
		if(r.readyState == 4){
			alert(r.response)
		}
	}
}

function cancelbuy() {
	var r = new XMLHttpRequest();
	r.open("DELETE", "/stocks/buy", true)
	r.setRequestHeader("Content-type", "application/json")
	r.send(JSON.stringify({TransactionNum: 1, UserId: user_id}))
	r.onreadystatechange = function(){
		if(r.readyState == 4){
			alert(r.response)
		}
	}
}

function sell() {
	stock = document.getElementById('sell_stock').value
	amount = document.getElementById('sell_amount').value
	var r = new XMLHttpRequest();
	r.open("POST", "/stocks/"+stock+"/sell", true)
	r.setRequestHeader("Content-type", "application/json")
	r.send(JSON.stringify({TransactionNum: 1, UserId: user_id, StockSymbol: stock, Amount: parseInt(amount)}))
	r.onreadystatechange = function(){
		if(r.readyState == 4){
			alert(r.response)
		}
	}
}

function commitsell() {
	var r = new XMLHttpRequest();
	r.open("PUT", "/stocks/sell", true)
	r.setRequestHeader("Content-type", "application/json")
	r.send(JSON.stringify({TransactionNum: 1, UserId: user_id}))
	r.onreadystatechange = function(){
		if(r.readyState == 4){
			alert(r.response)
		}
	}
}

function cancelsell() {
	var r = new XMLHttpRequest();
	r.open("DELETE", "/stocks/sell", true)
	r.setRequestHeader("Content-type", "application/json")
	r.send(JSON.stringify({TransactionNum: 1, UserId: user_id}))
	r.onreadystatechange = function(){
		if(r.readyState == 4){
			alert(r.response)
		}
	}
}

function setbuy() {
	stock = document.getElementById('set_buy_stock').value
	amount = document.getElementById('set_buy_amount').value
	var r = new XMLHttpRequest();
	r.open("POST", "/triggers/"+stock+"/buy", true)
	r.setRequestHeader("Content-type", "application/json")
	r.send(JSON.stringify({TransactionNum: 1, UserId: user_id, StockSymbol: stock, Amount: parseInt(amount)}))
	r.onreadystatechange = function(){
		if(r.readyState == 4){
			alert(r.response)
		}
	}
}

function setbuycancel() {
	stock = document.getElementById('set_buy_cancel_stock').value
	var r = new XMLHttpRequest();
	r.open("DELETE", "/triggers/"+stock+"/buy", true)
	r.setRequestHeader("Content-type", "application/json")
	r.send(JSON.stringify({TransactionNum: 1, UserId: user_id, StockSymbol: stock}))
	r.onreadystatechange = function(){
		if(r.readyState == 4){
			alert(r.response)
		}
	}
}

function setbuytrigger() {
	stock = document.getElementById('set_buy_trigger_stock').value
	amount = document.getElementById('set_buy_trigger_amount').value
	var r = new XMLHttpRequest();
	r.open("PUT", "/triggers/"+stock+"/buy", true)
	r.setRequestHeader("Content-type", "application/json")
	r.send(JSON.stringify({TransactionNum: 1, UserId: user_id, StockSymbol: stock, Amount: parseInt(amount)}))
	r.onreadystatechange = function(){
		if(r.readyState == 4){
			alert(r.response)
		}
	}
}

function setsell() {
	stock = document.getElementById('set_sell_stock').value
	amount = document.getElementById('set_sell_amount').value
	var r = new XMLHttpRequest();
	r.open("POST", "/triggers/"+stock+"/sell", true)
	r.setRequestHeader("Content-type", "application/json")
	r.send(JSON.stringify({TransactionNum: 1, UserId: user_id, StockSymbol: stock, Amount: parseInt(amount)}))
	r.onreadystatechange = function(){
		if(r.readyState == 4){
			alert(r.response)
		}
	}
}

function setsellcancel() {
	stock = document.getElementById('set_sell_cancel_stock').value
	var r = new XMLHttpRequest();
	r.open("DELETE", "/triggers/"+stock+"/sell", true)
	r.setRequestHeader("Content-type", "application/json")
	r.send(JSON.stringify({TransactionNum: 1, UserId: user_id, StockSymbol: stock}))
	r.onreadystatechange = function(){
		if(r.readyState == 4){
			alert(r.response)
		}
	}
}

function setselltrigger() {
	stock = document.getElementById('set_sell_trigger_stock').value
	amount = document.getElementById('set_sell_trigger_amount').value
	var r = new XMLHttpRequest();
	r.open("PUT", "/triggers/"+stock+"/sell", true)
	r.setRequestHeader("Content-type", "application/json")
	r.send(JSON.stringify({TransactionNum: 1, UserId: user_id, StockSymbol: stock, Amount: parseInt(amount)}))
	r.onreadystatechange = function(){
		if(r.readyState == 4){
			alert(r.response)
		}
	}
}

function dumplog() {
	user = document.getElementById('dumplog_user').value
	file = document.getElementById('dumplog_file').value
	var r = new XMLHttpRequest();
	r.open("POST", "/users/dump", true)
	r.setRequestHeader("Content-type", "application/json")
	r.send(JSON.stringify({TransactionNum: 1, UserId: user, FileName: file}))
	r.onreadystatechange = function(){
		if(r.readyState == 4){
			alert(r.response)
		}
	}
}

function dumplogall() {
	file = document.getElementById('dumplog_all_file').value
	var r = new XMLHttpRequest();
	r.open("POST", "/dump", true)
	r.setRequestHeader("Content-type", "application/json")
	r.send(JSON.stringify({TransactionNum: 1, FileName: file}))
	r.onreadystatechange = function(){
		if(r.readyState == 4){
			alert(r.response)
		}
	}
}

function summary() {
	user = document.getElementById('summary_user').value
	var r = new XMLHttpRequest();
	r.open("POST", "/users/"+user+"/summary", true)
	r.setRequestHeader("Content-type", "application/json")
	r.send(JSON.stringify({TransactionNum: 1, UserId: user}))
	r.onreadystatechange = function(){
		if(r.readyState == 4){
			alert(r.response)
		}
	}
}