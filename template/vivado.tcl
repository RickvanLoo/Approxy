link_design -part {{.Settings.PartName}}
read_vhdl [glob *.vhd]
synth_design {{if .Settings.OOC}}-mode out_of_context{{end}} {{if .Settings.NO_DSP}}-max_dsp 0{{end}} -top {{.TopName}}
create_clock -period 10.000 -name CLK -waveform {0.000 5.000} [get_ports {clk}]
{{if .Settings.WriteCheckpoint}}write_checkpoint -force {{.TopName}}_postsynth.dcp{{end}}
{{if .Settings.Placement}}place_design{{end}}
{{if .Settings.Route}}route_design{{end}}
{{if .Settings.Funcsim}}write_vhdl -mode funcsim {{.TopName}}_funcsim.vhd{{end}}
{{if and .Settings.Placement .Settings.WriteCheckpoint}}write_checkpoint -force {{.TopName}}_postplace.dcp{{end}}
{{if .Settings.Utilization}}report_utilization {{if .Settings.Hierarchical}}-hierarchical{{end}} -file {{.TopName}}_post_place_ult.rpt{{end}}
close_project
