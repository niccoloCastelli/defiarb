package price_graph

import (
	"github.com/niccoloCastelli/defiarb/storage/models"
	"sort"
	"strings"
)

const (
	directTreshold     = 1.014
	triangularTreshold = 1.015
)

func NewPriceGraph(pools ...models.LiquidityPool) *PriceGraph {
	poolsMap := map[string]models.LiquidityPool{}
	tokensMap := map[string]models.Token{}
	for i, pool := range pools {
		if _, ok := tokensMap[pool.Token0Address]; !ok {
			tokensMap[pool.Token0Address] = pool.Token0
		}
		if _, ok := tokensMap[pool.Token1Address]; !ok {
			tokensMap[pool.Token1Address] = pool.Token1
		}
		poolsMap[pool.Address] = pools[i]
	}
	//gr := graphs.GetGraph()
	return (&PriceGraph{
		pools:  poolsMap,
		tokens: tokensMap,
		graph:  NewGraph(),
	}).CalcGraph()
}

type PriceGraph struct {
	pools  map[string]models.LiquidityPool
	tokens map[string]models.Token
	graph  *Graph
}

func (g *PriceGraph) GetPrice(address string) (float64, float64, bool) {
	return 0, 0, false
}
func (g *PriceGraph) UpdatePrice(address string, value0 float64, value1 float64) *PriceGraph {
	if lp, ok := g.pools[address]; ok {
		_ = g.graph.SetWeight(EdgeKey(address, lp.Token0Address, lp.Token1Address), value0)
		_ = g.graph.SetWeight(EdgeKey(address, lp.Token1Address, lp.Token0Address), value1)
	}
	return g
}
func (g *PriceGraph) CalcGraph() *PriceGraph {
	g.graph = NewGraph()
	for addr, token := range g.tokens {
		g.graph.AddNode(*NewNode(addr, token.Symbol))
	}
	for addr, pool := range g.pools {
		if err := g.graph.AddEdge(pool.Token0Address, pool.Token1Address, pool.Token0Price, addr, pool.Description, pool.Exchange); err != nil {
			return nil
		}
		if err := g.graph.AddEdge(pool.Token1Address, pool.Token0Address, pool.Token1Price, addr, pool.Description, pool.Exchange); err != nil {
			return nil
		}
	}
	return g
}
func (g PriceGraph) pathKey(edges ...*Edge) string {
	addresses := make([]string, len(edges))
	for i, edge := range edges {
		addresses[i] = edge.Address
	}
	sort.Strings(addresses)
	return strings.Join(addresses, "/")
}
func (g *PriceGraph) ShortestPath(address string) ([]ArbPath, error) {
	opportunities := map[string]ArbPath{}
	node, err := g.graph.GetNode(address)
	if err != nil {
		return nil, err
	}
	for _, edge := range node.edgesFrom {
		weight := edge.Weight
		for _, e := range edge.To.edgesFrom {
			eWeight := weight * e.Weight
			if e.To.Address == node.Address {
				if edge.Address != e.Address && eWeight > directTreshold {
					pathKey := g.pathKey(edge, e)
					//fmt.Println("direct opportunity", edge.Name, e.Name, eWeight)
					arbPath := ArbPath{
						StartAddress:  edge.Address,
						StartExchange: edge.Exchange,
						Key:           g.pathKey(edge, e),
						Path:          []string{edge.Address, e.Address},
						Prices:        []float64{edge.Weight, e.Weight},
						Tokens:        []string{edge.From.Address, edge.To.Address},
						Exchanges:     []string{edge.Exchange, e.Exchange},
						Weight:        eWeight,
					}
					arbPath.Name = strings.Join([]string{edge.From.Name, edge.To.Name}, "->")
					opportunities[pathKey] = arbPath
				} else {
					continue
				}
			}
			for _, e2 := range e.To.edgesFrom {
				if e2.To.Address == node.Address {
					if eWeight*e2.Weight > triangularTreshold {
						edges := []*Edge{edge, e, e2}
						pathKey := g.pathKey(edge, e, e2)
						tokenNames := make([]string, len(edges))
						arbPath := ArbPath{
							StartAddress:  edges[0].Address,
							StartExchange: edges[0].Exchange,
							Key:           pathKey,
							Path:          make([]string, len(edges)),
							Prices:        make([]float64, len(edges)),
							Tokens:        make([]string, len(edges)),
							Exchanges:     make([]string, len(edges)),
							Weight:        edge.Weight * e.Weight * e2.Weight,
						}
						for i, edge := range edges {
							arbPath.Path[i] = edge.Address
							arbPath.Prices[i] = edge.Weight
							arbPath.Tokens[i] = edge.From.Address
							arbPath.Exchanges[i] = edge.Exchange
							tokenNames[i] = edge.From.Name
						}
						arbPath.Name = strings.Join(tokenNames, "->")
						opportunities[pathKey] = arbPath
					}
				}
			}
		}
	}

	ret := make([]ArbPath, 0, len(opportunities))
	for _, p := range opportunities {
		ret = append(ret, p)
		// fmt.Println("Opportunity!!!", edges[0].From.Name, edges[1].From.Name, edges[2].From.Name, edges[0].Weight*edges[1].Weight*edges[2].Weight, edges[0].Weight, edges[1].Weight, edges[2].Weight, time.Now())
	}
	return ret, nil

}

type ArbPath struct {
	StartExchange string
	StartAddress  string
	Key           string
	Name          string
	Weight        float64
	Path          []string
	Prices        []float64
	Tokens        []string
	Exchanges     []string
}
