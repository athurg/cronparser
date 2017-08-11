// vim: ft=html
package main

var htmlTemplate = `
<html>
	<head>
		<link rel="stylesheet" href="https://cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
		<script src="https://cdn.bootcss.com/jquery/1.12.4/jquery.min.js"></script>
		<script src="https://cdn.bootcss.com/bootstrap/3.3.7/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa" crossorigin="anonymous"></script>

		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bootstrap-datepicker/1.7.1/css/bootstrap-datepicker.min.css" integrity="sha256-I4gvabvvRivuPAYFqevVhZl88+vNf2NksupoBxMQi04=" crossorigin="anonymous" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/bootstrap-datepicker/1.7.1/js/bootstrap-datepicker.min.js" integrity="sha256-TueWqYu0G+lYIimeIcMI8x1m14QH/DQVt4s9m/uuhPw=" crossorigin="anonymous"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/bootstrap-datepicker/1.7.1/locales/bootstrap-datepicker.zh-CN.min.js" integrity="sha256-TWeu6bCkUfGPvdsbuXV9bNClNqZBR/WVdA6K8KFVSQA=" crossorigin="anonymous"></script>

		<style>
			table#result td{padding:0px;}
			table#result tr{padding:0px;}
			table#result th{padding:0px;}
			.tooltip-inner {max-width: 800px;}
		</style>
	</head>
	<body>
	<div class="container-fluid">
		<div class="page-header"><h1>Cron schedule view</h1></div>

		<div class="row">
			<form class="col-md-6 form-inline">
				<div class="form-group">
					<label for="start">起始日期</label>
					<input type="text" class="input-sm input-date form-control" id='start' name="start" value="{{.start.Format "2006-01-02"}}" />
				</div>
				<div class="form-group">
					<label for="end">结束日期</label>
					<input type="text" class="input-sm input-date form-control" id='end' name="end" value="{{.end.Format "2006-01-02"}}" />
				</div>
				<div class="checkbox">
					<label><input type="checkbox" name='showall' {{with .showall}}checked="checked"{{end}}>日期不限</label>
				</div>
				<button type="submit" class="btn btn-default">刷新</button>
			</form>
		</div>
		<table class="table table-bordered table-hover table-condensed" id="result">
			<thead>
				<tr>
					<th>分钟=></th>
					{{range $k,$v := (index .scheduleMap 0) }}
					<td>{{printf "%02d" $k}}</td>
					{{end}}
				</tr>
			</thead>
			<tbody>
				{{range $hour, $hourMap := .scheduleMap}}
				<tr>
					<th>{{printf "%02d:00" $hour}}</th>
					{{range $hourMap}}
					<td>
						{{with .}}
							<span class='badge' data-toggle='tooltip' data-placement='bottom' title='{{range .}}{{.}}<br/>{{end}}'>
								{{len .}}
							</span>
						{{end}}
						</td>
					{{end}}
				</tr>
				{{end}}
			</tbody>
		</table>

		<div class='row'>
			<ul class='list-group'>
				{{range .commands}}
				<li class='list-group-item'>{{.}}</li>
				{{end}}
			</ul>
		</div>
	</div>
	</body>
	<script>
	$('.input-date').datepicker({todayHighlight: true,todayBtn:"linked",format:'yyyy-mm-dd',language:'zh-CN'});
	$('[data-toggle="tooltip"]').tooltip({html:true, tooltipClass: "tooltip-style-fix" });
	$('input[name=showall]').change(function(){
		$('.input-date').attr("disabled", $('input[name=showall]').prop("checked"));
	});
	$('input[name=showall]').change();
	</script>

</html>
`
