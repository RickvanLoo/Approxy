open_checkpoint {{.TopName}}_postplace.dcp
read_saif {{.TopName}}_dump.saif
report_power -format xml -file {{.TopName}}_post_place_power.rpt
close_project
