link_design -part {{.Settings.PartName}}
read_vhdl [glob *.vhd]
synth_design {{if .Settings.OOC}}-mode out_of_context{{end}} {{if .Settings.NODSP}}-max_dsp 0{{end}} -top {{.TopName}}
{{if .Settings.Clk}}create_clock -period 10.000 -name CLK -waveform {0.000 5.000} [get_ports {clk}]{{end}}
{{if .Settings.WriteCheckpoint}}write_checkpoint -force {{.TopName}}_postsynth.dcp{{end}}
{{if .Settings.Placement}}place_design{{end}}
{{if .Settings.Route}}route_design{{end}}
{{if .Settings.Funcsim}}write_vhdl -mode funcsim {{.TopName}}_funcsim.vhd{{end}}
{{if and .Settings.Placement .Settings.WriteCheckpoint}}write_checkpoint -force {{.TopName}}_postplace.dcp{{end}}
{{if .Settings.Utilization}}report_utilization {{if .Settings.Hierarchical}}-hierarchical{{end}} -file {{.TopName}}_post_place_ult.rpt
set fo [open {{.TopName}}_primitive.rpt a]
puts $fo [llength [get_cells -hier -filter {PRIMITIVE_GROUP == CARRY}]]
close $fo
{{end}}
{{if .Settings.Timing}}report_timing -nworst 1 -path_type end -file {{.TopName}}_post_place_time.rpt{{end}}

close_project
