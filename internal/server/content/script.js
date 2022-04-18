/*

scenario = {
    msec_min,
    msec_max,
    workers,  // with instance
    tasks = {
        task_id: status,
        ...
    },
    items = [
        {
            worker_id
            worker_inst
            task_id
            msec_start
            msec_end
            status

            tooltip: function ()
            worker: function ()
        },
        ...
    ]
}

*/
function onlyUnique(value, index, self) {
    return self.indexOf(value) === index;
}

function demoJson2Scenario(jsonData) {
    var items = jsonData.map(function (d) {
        return {
            task_id: d.task_id,
            status: d.status,
            worker_id: d.worker_id,
            worker_inst: d.worker_inst,
            msec_start: (new Date(d.time_start)).getTime(),
            msec_end: (new Date(d.time_end)).getTime(),
            worker: function () {
                return `${this.worker_id}[${this.worker_inst}]`
            },

            tooltip: function () {
                return "task: " + this.task_id +
                    "\nstart: " + this.msec_start +
                    "\nend: " + this.msec_end +
                    "\nresult: " + d.label;
            }
        }
    })

    var tasks = {}
    jsonData.forEach(d => {
        var tid = d.task_id;
        var status = tasks[tid];
        if (status === undefined) {
            tasks[tid] = d.status
        } else {
            switch (d.status) {
                case "success":
                    tasks[tid] = d.status;
                case "error":
                    if (status != "success") {
                        tasks[tid] = d.status;
                    }
            }
        }
    });

    // get min and max timestamp
    var msec_min = d3.min(items, d => d.msec_start);
    var msec_max = d3.max(items, d => d.msec_end);

    // shift timestamps to msec_min
    items = items.map(function (d) {
        d.msec_start -= msec_min;
        d.msec_end -= msec_min;
        return d;
    })

    return {
        msec_min: msec_min,
        msec_max: msec_max,
        workers: d3.sort(d3.map(items, d => d.worker()).filter(onlyUnique)),
        tasks: tasks,
        items: items,
    }
}


function showInfo(scenario) {
    var tasks_success = 0;
    Object.values(scenario.tasks).forEach(function (status) {
        if (status == "success") tasks_success++;
    })

    var tasks_error = 0;
    Object.values(scenario.tasks).forEach(function (status) {
        if (status == "error") tasks_error++;
    })

    var msecs = {}
    scenario.items.forEach(function (item) {
        var msec = msecs[item.task_id];
        if (msec === undefined) {
            msecs[item.task_id] = {
                min: item.msec_start,
                max: item.msec_end,
            }
        } else {
            if (msec.min > item.msec_start) {
                msec.min = item.msec_start
            }
            if (msec.max > item.msec_end) {
                msec.max = item.msec_end
            }
            msecs[item.task_id] = msec
        }
    })

    var task_elapsed_min = 9999999999;
    var task_elapsed_max = 0;
    var elapsed_tot = 0;
    Object.values(msecs).forEach(function (msec) {
        var elapsed = msec.max - msec.min;
        if (task_elapsed_min > elapsed) {
            task_elapsed_min = elapsed
        }
        if (task_elapsed_max < elapsed) {
            task_elapsed_max = elapsed
        }
        elapsed_tot += elapsed;
    })

    var tasks_tot = Object.keys(scenario.tasks).length;

    var info = [
        ["msec_elapsed", scenario.msec_max - scenario.msec_min],
        ["workers", scenario.workers.length],
        ["tasks", tasks_tot],
        ["success", tasks_success],
        ["error", tasks_error],
        ["task elapsed min", task_elapsed_min],
        ["task elapsed avg", elapsed_tot / tasks_tot],
        ["task elapsed max", task_elapsed_max]
    ];

    d3.select("#info")
        .selectAll("table")
        .remove()
        ;
    d3.select("#info")
        .append("table")
        .selectAll("tr")
        .data(info)
        .join("tr")
        .selectAll("td")
        .data(d => d)
        .join("td")
        .text(d => d);

}

function drawGraphWorkers(scenario) {
    var margin = { top: 10, right: 40, bottom: 30, left: 100 },
        width = 1600 - margin.left - margin.right,
        height = 800 - margin.top - margin.bottom;

    // remove old graph
    d3.select("#graphWorkers").select("svg").remove();

    // append the svg object to the body of the page
    var svG = d3.select("#graphWorkers")
        .append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .append("g")
        .attr("transform",
            "translate(" + margin.left + "," + margin.top + ")");

    // X scale and Axis
    var x = d3.scaleLinear()
        .domain([0, scenario.msec_max - scenario.msec_min])
        .range([0, width]);
    svG
        .append('g')
        .attr("transform", "translate(0," + height + ")")
        .call(d3.axisBottom(x));

    // Y scale and Axis
    var y = d3.scaleBand()
        .domain(scenario.workers)
        .range([height, 0])
        .padding([0.2]);
    svG
        .append('g')
        .call(d3.axisLeft(y));

    // draw items
    svG
        .selectAll("whatever")
        .data(scenario.items)
        .enter()
        .append("rect")
        .attr("x", d => x(d.msec_start))
        .attr("y", d => y(d.worker()))
        .attr("width", d => x(d.msec_end - d.msec_start))
        .attr("height", y.bandwidth())
        .attr("class", d => d.status)
        .on("click", fnClick)
        .append("title")
        .text(d => d.tooltip())
        ;

};

function drawGraphTasks(scenario) {
    var margin = { top: 10, right: 40, bottom: 30, left: 100 },
        width = 1600 - margin.left - margin.right,
        height = 800 - margin.top - margin.bottom;

    // remove old graph
    d3.select("#graphTasks").select("svg").remove();

    // append the svg object to the body of the page
    var svG = d3.select("#graphTasks")
        .append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .append("g")
        .attr("transform",
            "translate(" + margin.left + "," + margin.top + ")");



    taskids = d3.sort(Object.keys(scenario.tasks));


    // X scale and Axis
    var x = d3.scaleLinear()
        .domain([0, scenario.msec_max - scenario.msec_min])
        .range([0, width]);
    svG
        .append('g')
        .attr("transform", "translate(0," + height + ")")
        .call(d3.axisBottom(x));

    // Y scale and Axis
    var y = d3.scaleBand()
        .domain(taskids)
        .range([height, 0])
    // .padding([0.2]);
    svG
        .append('g')
        .call(d3.axisLeft(y)
            // .tickValues(["t010", "t020"])
        );

    // draw items
    svG
        .selectAll("whatever")
        .data(scenario.items)
        .enter()
        .append("rect")
        .attr("x", d => x(d.msec_start))
        .attr("y", d => y(d.task_id))
        .attr("width", d => x(d.msec_end - d.msec_start))
        .attr("height", y.bandwidth())
        .attr("class", d => d.status)
        .on("click", fnClick)
        .append("title")
        .text(d => d.tooltip())
        ;

};



function onchange() {
    selectValue = selectjsons.property('value');
    d3.json("/data/" + selectValue).then(
        function (jsonData) {
            var scenario = demoJson2Scenario(jsonData);
            drawGraphWorkers(scenario);
            drawGraphTasks(scenario);
            showInfo(scenario);
        }
    );
};

function fnClick(event, data) {
    var task_id = data.task_id;
    d3.select("#graphWorkers")
        .selectAll("rect")
        .nodes()
        .map(function (d) {
            var task = d3.select(d);
            task.classed("highlight", task.datum().task_id == task_id);
        })
        ;
    d3.select("#graphTasks")
        .selectAll("rect")
        .nodes()
        .map(function (d) {
            var task = d3.select(d);
            task.classed("highlight", task.datum().task_id == task_id);
        })
        ;
}


var selectjsons = d3.select("#filejson")


// bind select change event
selectjsons.on('change', onchange);

// fill the select with the data json files
d3.json("/data").then(
    function (data) {
        // console.log(data);
        selectjsons
            .selectAll('option')
            .data(data)
            .join("option")
            .text(d => d)
            .attr("value", d => d)
            ;
        onchange();
    }
)
