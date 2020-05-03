package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/Nosp1/TA/Covid19/analytics"
	"github.com/teamwork/reload"
)

var opened = false

var CSS string = `<style>
		html, body {
		width: 100%;
		min-height: 100%;
		overflow-x: hidden;
		margin: 0;
		-webkit-transition: all 0.2s ease 0s;
		transition: all 0.2s ease 0s;
	}
	body {
		font-family: 'roboto', sans-serif;
		font-size: 18px;
		font-weight:900;
		-webkit-font-smoothing: antialiased;
		color: #2c3e50;
		text-align: center;
	}

	main {
		position: relative;
		background-color: lavender;
		margin-right: auto;
		margin-left: auto;
		margin-top: 150px;
		margin-bottom: 50px;
		border-style: solid;
		border-color: white;
		border-radius: 5px;
		color: #333333;
		width: 800px;
		min-height: 300px;
		padding: 50px;
		-webkit-box-shadow: 0px 1px 15px -4px rgba(0,0,0,0.75);
		-moz-box-shadow: 0px 1px 15px -4px rgba(0,0,0,0.75);
		box-shadow: 0px 1px 15px -4px rgba(0,0,0,0.75);
	}

	h1 {
		margin-top: 20px;
	}

	table {
		border-collapse: collapse;
		width:200px;
	}

	th, td {
		text-align: center;
		padding: 8px;
		min-width: 70px;
		font-weight: 500;
	}

	td {
		font-weight: 400;
	}

tr:nth-child(even){
		background-color: #f2f2f2
	}
	p {
		margin: 15px 3px 12px 3px
	}

	#box {
		position: relative;
		width: 800px;
		min-height: 300px;
		background-color:blue;
		display: flex;
		flex-direction: row;
		flex-wrap: wrap;
		justify-content: center;
		align-items: center;
	}
	#kategori {
		height: 100%;
		width: 200px;
		font-weight: 500;
		margin: 5px;
		margin-top: 111px;
		margin-left: 40px;
		margin-bottom: 80px;
		text-align: right;
	}
	#tabell {
		height: 100%;
		width: 400px;
		margin: 5px;
		margin-top: 8px;
	}

	</style>`

var beginstyle string = `<main>
	<div id="box">
		<h1 style="width: 100%; margin-top: 50px;">Covid-19 today </h1>
<div id="kategori">
<p> infected total:</p>
<p>New infected today:</p>
<p>Death total:</p>
<p>new Deaths today:</p>
</div>
<div id="tabell">
<table>
<tr>
<th>Today</th>
<th>yesterday</th>
<th>Difference</th>
</tr>`

var endstyle string = `</table>
            </div>
        </div>
    </main>`

var stats = analytics.GetStats()
var queue = []string{
	stats[0].TotalInfected, stats[1].TotalInfected,
	stats[0].NewInfected, stats[1].NewInfected,
	stats[0].TotalDeath, stats[1].TotalDeath,
	stats[0].NewDeath, stats[1].NewDeath,
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, CSS)
	fmt.Fprintf(w, beginstyle)
	fmt.Fprintf(w, "<tr>")

	lineshift := 0

	for i := 0; i < 4; i++ {
		today, _ := strconv.Atoi(queue[lineshift])
		yesterday, _ := strconv.Atoi(queue[lineshift+1])
		difference := today - yesterday
		differconverter := math.Abs(float64(difference))

		fmt.Fprintf(w, "<tr>")
		for j := 0; j < 3; j++ {
			if j == 2 {
				if difference < 0 {
					fmt.Fprintf(w, "<td> "+strconv.Itoa(int(differconverter))+"</td>")
				} else {
					fmt.Fprintf(w, "<td> "+strconv.Itoa(int(differconverter))+"</td>")
				}
			} else {
				fmt.Fprintf(w, "<td>"+queue[j+lineshift]+"</td>")
			}
		}
		fmt.Fprintf(w, "</tr>")
		lineshift += 2

	}
	fmt.Fprintf(w, endstyle)
}

func main() {

	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
	timer := time.NewTimer(1 * time.Hour)

	<-timer.C
	reload.Exec()

}
