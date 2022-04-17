fieldname = {
    task_id: "isin",
    worker_id: "source",
    worker_inst: "instance",
    status: "status",
    label: "label",
    err: "err",
    time_start: "time_start",
    time_end: "time_end",
};



function quotesJson2Data(d, msec_min) {
    var msec_start = (new Date(d.time_start)).getTime() - msec_min;
    var msec_end = (new Date(d.time_end)).getTime() - msec_min;
    var msec_elapsed = msec_end - msec_start;

    var tooltip = "isin: " + d.isin +
        "\nstart: " + msec_start +
        "\nend: " + msec_end;

    var state;
    if (typeof d.error === 'undefined') {
        state = 'success';
        tooltip = tooltip + "\nresult: " + d.currency + " " + d.price;
    } else if (d.error.includes('context canceled')) {
        state = 'canceled';
    } else {
        state = 'error';
    }

    return {
        worker: d.source + "[" + d.instance + "]",
        task: d.isin,
        msec_start: msec_start,
        msec_end: msec_end,
        msec_elapsed: msec_elapsed,
        class: state,
        tooltip: tooltip
    }
}


function demoJson2Data(d, msec_min) {
    var msec_start = (new Date(d.time_start)).getTime() - msec_min;
    var msec_end = (new Date(d.time_end)).getTime() - msec_min;
    var msec_elapsed = msec_end - msec_start;

    var tooltip = "task: " + d.task_id +
        "\nstart: " + msec_start +
        "\nend: " + msec_end +
        "\nresult: " + d.label;

    return {
        worker: d.worker_id + "[" + d.worker_inst + "]",
        task: d.task_id,
        msec_start: msec_start,
        msec_end: msec_end,
        msec_elapsed: msec_elapsed,
        class: d.status,
        tooltip: tooltip
    }
}


function onlyUnique(value, index, self) {
    return self.indexOf(value) === index;
}



function drawGraph(jsonData) {
    var tmin = d3.min(jsonData, d => d[fieldname.time_start]);
    var tmax = d3.max(jsonData, d => d[fieldname.time_end]);

    var msec_min = (new Date(tmin)).getTime();
    var msec_max = (new Date(tmax)).getTime();
    var msec_elapsed = msec_max - msec_min;

    console.log("msec_min =", msec_min);
    console.log("msec_max =", msec_max);
    console.log("msec_elapsed =", msec_elapsed);

    // transform jsonData to data
    var data = d3.map(jsonData, function (d) {
        if (d.isin !== undefined) {
            return quotesJson2Data(d, msec_min);
        } else {
            return demoJson2Data(d, msec_min);
        }
    });

    // filter data items with small elapsed 
    // data = d3.filter(data, d => d.msec_elapsed >= 10);

    var workers = d3.sort(d3.map(data, d => d.worker).filter(onlyUnique));
    console.log(workers);

    var tasks = d3.sort(d3.map(data, d => d.task).filter(onlyUnique));
    console.log(tasks);

    var xleft = 150;
    var margin = { top: 10, right: 40, bottom: 30, left: xleft },
        width = 1600 - margin.left - margin.right,
        height = 800 - margin.top - margin.bottom;

    d3.select("#graph").select("svg").remove();

    // append the svg object to the body of the page
    var svG = d3.select("#graph")
        .append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .append("g")
        .attr("transform",
            "translate(" + margin.left + "," + margin.top + ")");

    // X scale and Axis
    var x = d3.scaleLinear()
        .domain([0, msec_elapsed])         // This is the min and the max of the data: 0 to 100 if percentages
        .range([0, width]);       // This is the corresponding value I want in Pixel
    svG
        .append('g')
        .attr("transform", "translate(0," + height + ")")
        .call(d3.axisBottom(x));

    // X scale and Axis
    var y = d3.scaleBand()
        .domain(workers)         // This is the min and the max of the data: 0 to 100 if percentages
        .range([height, 0])       // This is the corresponding value I want in Pixel
        .padding([0.2]);
    svG
        .append('g')
        .call(d3.axisLeft(y));

    // Add 3 dots for 0, 50 and 100%
    svG
        .selectAll("whatever")
        .data(data)
        .enter()
        .append("rect")
        .attr("x", d => x(d.msec_start))
        .attr("y", d => y(d.worker))
        .attr("width", d => x(d.msec_elapsed))
        .attr("height", y.bandwidth())
        .attr("class", d => d.class)
        //      .attr("alt", d => d.isin)
        //      .attr("fill", d => isinFillColor(d) )
        .on("click", fnClick)
        .append("title")
        .text(d => d.tooltip)
        ;


} // drawGraph


function fnClick(event, data) {
    // console.log(event.target);
    // event.target.classList.toggle("highlight");

    // console.log(data);

    var taskid = data.task;

    d3.select("#graph")
        .selectAll("rect")
        .nodes()
        .map(function (d) {
            //console.log(d);

            var task = d3.select(d);
            task.classed("highlight", task.datum().task == taskid);
            //d.classList.toggle("highlight");
            //console.log(task.datum());
        })
        ;


}


var selectjsons = d3.select("#filejson")

function onchange() {
    selectValue = selectjsons.property('value');
    d3.json("/data/" + selectValue).then(
        function (data) {
            drawGraph(data);
        }
    );
};


selectjsons.on('change', onchange);

// fill the select with the data json files
// d3.json("/data").then(
//     function (data) {
//         console.log(data);
//         selectjsons
//             .selectAll('option')
//             .remove()
//             .data(data)
//             .enter()
//             .append('option')
//             .text(d => d)
//             .attr("value", d => d)
//             ;
//         onchange();
//     }
// )

d3.json("/data").then(
    function (data) {
        console.log(data);
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
