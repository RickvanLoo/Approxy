# "lpACLib" is a library for Low-Power Approximate Computing Modules.
# Copyright (C) 2016, Walaa El-Harouni, Muhammad Shafique, CES, KIT.
# E-mail: walaa.elharouny@gmail.com, swahlah@yahoo.com

# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License
# along with this program. If not, see <http://www.gnu.org/licenses/>.
#========================================================================

sh mkdir -p synopsys
set objects synopsys

# code to easier output timestamped files later, next to non-timestamped files:
set initTime [clock seconds]

set TIMESTAMP [clock format $initTime -format {%Y-%m-%d-%H:%M}]

# defining clock
set CLK  "clk"
set RST  "reset"
set CLK_PERIOD 3.22; # 3.22 ns (310 MHz)
set CLK_UNCERTAINTY  0.1; # 200ps

# this function will run the specified command, and copy the file to a timestamped
# version, for multiple synthesize runs
# cmd_and_timecopy cmd filepath suffix
# with cmd="command including the output director (can be for example -o or can be ">")"
# filepath="filename to be written WITHOUT file extension"
# suffix="file extension, for example txt"

proc cmd_and_timecopy {cmd filepath suffix} {
    global TIMESTAMP
    eval $cmd $filepath.$suffix
    sh cp -f $filepath.$suffix ${filepath}_$TIMESTAMP.$suffix
}
# --------------------------------------------------------------------------

# there are sometimes issues with relative paths, so we have to use a BASE_DIR
# variable
set BASE_DIR /home/harouni/WORK/SAD #TODO: change path to work directory

#set SIM_DIR ./

set DESIGN_NAME SAD8x8Zero #TODO: change to top entity name

# generate a directory structure in which we will later write our reports, etc.
sh mkdir -p $BASE_DIR/syn/report
sh mkdir -p $BASE_DIR/syn/ddc
sh mkdir -p $BASE_DIR/syn/log

# generated by the leon3mp project with "make scripts"
#set trans_dc_max_depth 1
#set hdlin_seqmap_sync_search_depth 1
#set hdlin_nba_rewrite false
set hdlin_ff_always_sync_set_reset true
set hdlin_ff_always_async_set_reset false
#set hdlin_infer_complex_set_reset true
#set hdlin_translate_off_skip_text true
set suppress_errors VHDL-2285
#set hdlin_use_carry_in true

# MAYBE needed to preserve names for correct annotated gate level simulation
set power_preserve_rtl_hier_names true

#TODO: list all source files
analyze -f VHDL {./SAD8x8Zero.vhd ./SAD8x1Zero.vhd ./AdderIMPACTZeroApproxMultiBit.vhd ./AdderIMPACTZeroApproxOneBit.vhd ./AdderAccurateOneBit.vhd}

# this can be used in further refinement-runs, to re-read the previous synthesized gate-level:
#read_ddc $BASE_DIR/syn/ddc/leon3.ddc
#read_sdc $BASE_DIR/syn/ddc/constraints_dc.sdc

echo "### Elaboration ###"

elaborate $DESIGN_NAME

echo "### Synthesis ###"

# select the current toplevel of the design, this is the leon3mp entity/module
current_design $DESIGN_NAME

#read_ddc $BASE_DIR/syn/ddc/elab.ddc

# link the design
cmd_and_timecopy "link >" "$BASE_DIR/syn/report/link" txt

# output .saif file for potential forward annotation
#cmd_and_timecopy "rtl2saif -design leon3mp -output" "$BASE_DIR/syn/ddc/leon3mp_fwd" saif

# remove multiple instantations for synthesis, unique design for each cell
uniquify

# this will optimize more, but flatten the hierarchy completely:
#ungroup -flatten -force -all
#ungroup -flatten -all

# check design and write out reports:
cmd_and_timecopy "check_design -summary >" "$BASE_DIR/syn/report/check_design_summary" txt
cmd_and_timecopy "check_design >" "$BASE_DIR/syn/report/check_design" txt

# report on hierarchy:
cmd_and_timecopy "report_hierarchy >" "$BASE_DIR/syn/report/hierarchy" txt

# read to convert library for vhdl/verilog later
#read_lib /home/sun/LEON3/library/TSMCHOME/digital/Front_End/timing_power_noise/NLDM/tcbn45gsbwp_120a/tcbn45gsbwpwc.lib
#quit

# ------
# BEGIN: DERIVED FROM timing.tcl and easic_timing.tcl:

set_wire_load_mode segmented
set auto_wire_load_selection "true"
set_wire_load_mode segmented

#set_critical_range 1.0 leon3mp

# END: DERIVED FROM timing.tcl and easic_timing.tcl
# ------

set_max_area 0

# report all constraint violations:
cmd_and_timecopy "report_constraints -all_violators >" "$BASE_DIR/syn/report/all_constraints" txt


create_clock $CLK -period $CLK_PERIOD
set_clock_uncertainty $CLK_UNCERTAINTY [all_clocks]
set_dont_touch_network [all_clocks]

remove_driving_cell $RST
set_drive 0 $RST
set_dont_touch_network $RST

#set_output_delay $DFF_SETUP -clock $CLK [all_outputs]
#set_load 0.2 [all_outputs]

#set all_inputs_wo_rst_clk [remove_from_collection [remove_from_collection [all_inputs] [get_port $CLK]] [get_port $RST]]
#set_input_delay -clock $CLK $DFF_CKQ $all_inputs_wo_rst_clk
#set_driving_cell -library $LIB_NAME -lib_cell $DFF_CELL -pin Q $all_inputs_wo_rst_clk

# Fix hold violations automatically
# http://www.tkt.cs.tut.fi/tools/public/tutorials/synopsys/design_compiler/gsdc.html
set_fix_hold $CLK

# Different options for the compile_ultra command (high effort optimization)
# compile_ultra is recommended for 90nm and smaller:

#compile_ultra -retime -area_high_effort_script -timing_high_effort_script -no_autoungroup -self_gating
#compile_ultra -no_autoungroup -gate_clock
#compile_ultra -retime
compile_ultra -exact_map -no_autoungroup

# report timing and area of synthesized design
cmd_and_timecopy "report_timing >" "$BASE_DIR/syn/report/timing_ultra" txt
cmd_and_timecopy "report_area >" "$BASE_DIR/syn/report/area_ultra" txt
cmd_and_timecopy "report_area -hierarchy >" "$BASE_DIR/syn/report/area_hierarchy" txt
cmd_and_timecopy "report_power >" "$BASE_DIR/syn/report/power_ultra" txt
cmd_and_timecopy "report_clock >" "$BASE_DIR/syn/report/clock" txt

echo "### Writing Gate Level files .. ###"

# vhdl/neutral-based synopsys ddc ##
cmd_and_timecopy "write -f ddc -hierarchy $DESIGN_NAME -o" "$BASE_DIR/syn/ddc/$DESIGN_NAME" ddc

# write out constraint informations in synopsys format:
cmd_and_timecopy "write_sdc" "$BASE_DIR/syn/ddc/constraints_dc" sdc

#### do not use change_names first to write vhd and verilog files


# write out constraint informations in standard delay format:
cmd_and_timecopy "write_sdf -version 2.1 -context vhdl" "$BASE_DIR/syn/ddc/constraints_dc_vhd1" sdf
#cmd_and_timecopy "write_sdf -version 2.1 -context verilog" "$BASE_DIR/syn/ddc/constraints_dc_vlog1" sdf

# now verilog based naming by change_names ###############################################
#echo "Verilog.."

# change names of gate level according to verilog-friendly names for the rest
#change_names -rules verilog -hierarchy

#cmd_and_timecopy "write_lib saed90nm_typ -format verilog -output" "gatelevel_lib_saed90nm_typ_vlog2" vl

# write out gate level design in synopsys proprietary ddc file format:
#cmd_and_timecopy "write -f ddc -hierarchy $DESIGN_NAME -o" "$BASE_DIR/syn/ddc/$DESIGN_NAME_vlog2" ddc

# write out gate level design in verilog format using standard cells from our library:
#cmd_and_timecopy "write -f verilog -hierarchy $DESIGN_NAME -o" "$BASE_DIR/syn/ddc/$DESIGN_NAME-gatelevel_vlog2" vl

# write out constraint informations in standard delay format:
#cmd_and_timecopy "write_sdf -version 2.1 -context verilog" "$BASE_DIR/syn/ddc/constraints_dc_vlog2" sdf

# vhdl based naming by change_names ###############################################
# not sure if this really should/does work like expecte, because the names might now me messed up, because
# change_names was already used with -rules verilog above!
echo "VHDL.."

# the same again with names changed according to vhdl rules explicitly
change_names -rules vhdl -hierarchy

#cmd_and_timecopy "write_lib saed90nm_typ -format vhdl -output" "gatelevel_lib_saed90nm_typ_vhd2" vhd
cmd_and_timecopy "write -f vhdl -hierarchy $DESIGN_NAME -o" "$BASE_DIR/syn/ddc/$DESIGN_NAME-gatelevel_vhd2" vhd
cmd_and_timecopy "write_sdf -version 2.1 -context vhdl" "$BASE_DIR/syn/ddc/constraints_dc_vhd2" sdf

quit
