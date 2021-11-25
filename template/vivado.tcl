link_design -part {{.PartName}}
read_vhdl [glob *.vhd]
synth_design {{if .OOC}}-mode out_of_context{{end}} -top {{.TopName}}
{{if .WriteCheckpoint}}write_checkpoint -force {{.TopName}}_postsynth.dcp{{end}}
{{if .Placement}}place_design{{end}}
{{if and .Placement .WriteCheckpoint}}write_checkpoint -force {{.TopName}}_postplace.dcp{{end}}
{{if .Utilization}}report_utilization -file {{.TopName}}_post_place_ult.rpt{{end}}
close_project
