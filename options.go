package WoComApi

type WoComOption func(com *WoCom)

func CorpId(id string) WoComOption {
	return func(com *WoCom) {
		com.corpId = id
	}
}

func Secret(secret string) WoComOption {
	return func(com *WoCom) {
		com.secret = secret
	}
}

func AgentId(agentid string) WoComOption {
	return func(com *WoCom) {
		com.agentId = agentid
	}
}
