package dao

import "fwallet/model"

func (d *Dao) ListNetwork(page, size int) (int64, []model.Network, error) {
	var beans []model.Network
	total, err := d.orm.Where("1=1").OrderBy("ID ASC").
		Limit(size, (page-1)*size).FindAndCount(&beans)
	return total, beans, err
}
