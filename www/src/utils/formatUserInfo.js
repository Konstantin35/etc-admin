import formatEtc from './formatEtc'

export default function(data){
	var users = []
	data.forEach(el => {
		var user = {
			name: el.BasicInfo.account,
			tel: el.BasicInfo.phone,
			email: el.BasicInfo.email,
			vip: el.BasicInfo.vip,
			wallet: []
		}
		var wallet = {
			address: el.BasicInfo.walletAddress,
			fee: el.BasicInfo.fee,
			lastBanefit: formatEtc(el.LastRevenue),
			totalBanefit: formatEtc(el.AllRevenue),
			stats: !!el.OfflineTime ? '离线' : '在线',
			offLineTime: el.offLineTime
		}

		var index = 0

		if(!users.some( (el,i) => {
			if(el.name === user.name){
				users[i].wallet.push(wallet)
				return true
			}
			index=i
			return false
		})){
		  users.push(user)
		  users[index].wallet.push(wallet)
		}
	})

  return users
}