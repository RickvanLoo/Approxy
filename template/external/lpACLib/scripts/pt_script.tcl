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

# code to easier output timestamped files later, next to non-timestamped files:
set initTime [clock seconds]

set TIMESTAMP [clock format $initTime -format {%Y-%m-%d-%H:%M}]

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
set BASE_DIR /home/harouni/WORK/2x2/accurateMult/modelsim/ #TODO: change to work directory

sh mkdir -p $BASE_DIR/report

set stdcells_home /home/harouni/TSMC_library/ #TODO: change to technology library path
set search_path "$stdcells_home"
set target_library "tcbn45gsbwpwc.db"
set link_path "* $target_library"
set power_enable_analysis "true"
read_vhdl "./Accurate2x2Mult-gatelevel_vhd2.vhd" #TODO: change to net-list file name
current_design "Accurate2x2Mult" #TODO: change to top entity
link
set power_analysis_mode "averaged"
#TODO: change to desired saif file name
read_saif "./accurate.saif" -strip_path "testbench/uut"
cmd_and_timecopy "report_switching_activity -list_not_annotated >" "$BASE_DIR/report/report_switching_activity_averaged" txt
cmd_and_timecopy "report_power -verbose -hierarchy >" "$BASE_DIR/report/report_power_averaged" txt
set power_analysis_mode "time_based"
#TODO: change to desired vcd file name
read_vcd "./accurate.vcd" -strip_path "testbench/uut"
cmd_and_timecopy "report_switching_activity -list_not_annotated >" "$BASE_DIR/report/report_switching_activity_time_based" txt
cmd_and_timecopy "report_power -verbose -hierarchy >" "$BASE_DIR/report/report_power_time_based" txt
quit
