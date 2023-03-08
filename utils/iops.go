package IOPS

type IOPS struct {
	Value     int `json:"value"`
	Blocksize int `json:"blocksize"`
}

func ConvertBlocksize(iops IOPS, toBlocksize int32) IOPS {
	return IOPS{Value: (iops.Value * iops.Blocksize / toBlocksize), Blocksize: toBlocksize}
}
