link_design -part {{.Settings.PartName}}
read_vhdl [glob *.vhd]
synth_design {{if .Settings.OOC}}-mode out_of_context{{end}} {{if .Settings.NO_DSP}}-max_dsp 0{{end}} -top {{.TopName}}
{{if .Settings.WriteCheckpoint}}write_checkpoint -force {{.TopName}}_postsynth.dcp{{end}}
{{if .Settings.Placement}}place_design{{end}}
{{if and .Settings.Placement .Settings.WriteCheckpoint}}write_checkpoint -force {{.TopName}}_postplace.dcp{{end}}
{{if .Settings.Utilization}}report_utilization -hierarchical -file {{.TopName}}_post_place_ult.rpt{{end}}
close_project
