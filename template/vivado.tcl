link_design -part {{.PartName}}
read_vhdl [glob {{.Folder}}/*.vhd]
synth_design {{if .OOC}}-mode out_of_context{{end}} -top {{.Top}}
{{if .WriteCheckpoint}}write_checkpoint -force {{.Folder}}/postsynth.dcp{{end}}
place_design
{{if .WriteCheckpoint}}write_checkpoint -force {{.Folder}}/postplace.dcp{{end}}
{{if .Utilization}}report_utilization -file {{.Folder}}/post_place_ult.rpt{{end}}
close_project
