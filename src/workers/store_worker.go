package workers

type StoreWorkerParam struct {
	CompanyID string
	ShopID    string
	PersonID  string
	CaptureAt string
}

func storeFrequentCustomerHandler(companyID string, shopID string, personID string, captureAt string) {
	// 来了个新客

	// 1. 首先看这组companyID shopID里有没有这个personID的bitmap，bitmap里记录了一个值，当天这人有没有来过

	// 1.1 有的话，这就是一个来过的人，记在bitmap中更新那一行

	// 1.2 没有的话，就没来过，给bitmap添加一行

	// 2. 根据bitmap中的频次，记到当天的数据分布表中, 总人数，高频次数，低频次数，新客数，总到访间隔天数，总到访天数

	// 3. 取一下频率规则，判断是不是高频次的人

	// 3.1 是的话，根据来的captureAt时间，记到高频表里

}
