-- "lpACLib" is a library for Low-Power Approximate Computing Modules.
-- Copyright (C) 2016, Walaa El-Harouni, Muhammad Shafique, CES, KIT.
—- E-mail: walaa.elharouny@gmail.com, swahlah@yahoo.com

-- This program is free software: you can redistribute it and/or modify
-- it under the terms of the GNU General Public License as published by
-- the Free Software Foundation, either version 3 of the License, or
-- (at your option) any later version.

-- This program is distributed in the hope that it will be useful,
-- but WITHOUT ANY WARRANTY; without even the implied warranty of
-- MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
-- GNU General Public License for more details.

-- You should have received a copy of the GNU General Public License
-- along with this program. If not, see <http://www.gnu.org/licenses/>.

--------------------------------------------
-- Testbench for all 32x32 SADs (all adder types and LSBs)
-- Author: Walaa El-Harouni  
--------------------------------------------

LIBRARY ieee;
USE ieee.std_logic_1164.ALL;

-- Define an array of 8 bytes
--package TypesDefinition is
--      type BYTE_ARRAY_8 is array(7 downto 0) of std_logic_vector(7 downto 0);
--      type BYTE_ARRAY_32 is array(3 downto 0) of BYTE_ARRAY_8;
--end TypesDefinition ;
--==================================================


LIBRARY ieee;
USE ieee.std_logic_1164.ALL;

library VITAL;
use VITAL.all;

ENTITY testbench IS
  generic ( bitWidth : integer := 22);
  port (
	outReady : out std_logic;
	SAD : out std_logic_vector(bitWidth-1 downto 0)
	);
END testbench;

use work.TypesDefinition.all;

ARCHITECTURE behavior OF testbench IS
	
	component SAD32x32First is
	generic (bitWidth : integer := 22; approxBits : integer := 0);
	port (A32   	: in BYTE_ARRAY_32;
	      B32 	: in BYTE_ARRAY_32;
	      i_valid	: in std_logic;
	      clk	: in std_logic;
	      reset	: in std_logic;
	      outReady  : out std_logic;
	      SAD 	: out std_logic_vector(bitWidth-1 downto 0));
	end component;

	signal a: BYTE_ARRAY_32;
	signal b: BYTE_ARRAY_32;
	--signal sad: std_logic_vector(21 downto 0);
	signal clock: std_logic;
	signal i_valid: std_logic;
	signal reset: std_logic;
	--signal outReady: std_logic;
	
	
	-- A { 0, 42, 6, 0, 58, 14 ,0 , 0 }
	-- B{ 0, 0, 16, 0, 96, 16, 0 , 0 }
	-- SAD should be 92
	
	-- A { 255, 5, 5, 4, 0, 7, 7, 5 }
	-- B { 255, 7, 5, 1, 7, 4, 8, 4 }
	-- SAD should be 17
	
	-- A { 234, 14, 255, 123, 134, 141, 201, 221 }
	-- B { 134, 141, 255, 117, 7, 141, 8, 222 }
	-- SAD should be 554
	
	-- A { 0, 42, 6, 0, 58, 14 ,0 , 0 }
	-- B { 255, 7, 5, 1, 7, 4, 8, 4 }
	-- SAD should be 365

BEGIN
	uut: SAD32x32First port map (A32 => a, B32 => b, i_valid => i_valid, clk => clock, reset => reset, outReady => outReady, SAD => sad);
		
	clk_proc : process
	begin
		clock <= '0';
		wait  for 1 ns;
		clock <= '1';
		wait for 1 ns;
	end process;
	
	--valid_proc : process
	--begin
	--	i_valid <= '1'
	--	wait  for 1 ns;
	--	i_valid <= '1';
	--	wait for 1 ns;
	--end process;
	
	stim_proc : process
	begin
		reset <= '1';
		i_valid <= '0';
		wait for 2 ns;
		reset <= '0';
		i_valid <= '1';
		wait for 1 ns;
		
		for l in 0 to 10000 loop
		a(0)(0) <= "00000000"; a(0)(1) <= "00101010"; a(0)(2) <= "00000110"; a(0)(3) <= "00000000";
		a(0)(4) <= "00111010"; a(0)(5) <= "00001110"; a(0)(6) <= "00000000"; a(0)(7) <= "00000000";

		b(0)(0) <= "00000000"; b(0)(1) <= "00000000"; b(0)(2) <= "00010000"; b(0)(3) <= "00000000";
		b(0)(4) <= "01100000"; b(0)(5) <= "00010000"; b(0)(6) <= "00000000"; b(0)(7) <= "00000000"; 
		
		-- ========================================================================================
		
		a(1)(0) <= "11111111"; a(1)(1) <= "00000101"; a(1)(2) <= "00000101"; a(1)(3) <= "00000100";
		a(1)(4) <= "00000000"; a(1)(5) <= "00000111"; a(1)(6) <= "00000111"; a(1)(7) <= "00000101";

		b(1)(0) <= "11111111"; b(1)(1) <= "00000111"; b(1)(2) <= "00000101"; b(1)(3) <= "00000001";
		b(1)(4) <= "00000111"; b(1)(5) <= "00000100"; b(1)(6) <= "00001000"; b(1)(7) <= "00000100"; 
		
		-- ========================================================================================
		
		a(2)(0) <= "11101010"; a(2)(1) <= "00001110"; a(2)(2) <= "11111111"; a(2)(3) <= "01111011";
		a(2)(4) <= "10000110"; a(2)(5) <= "10001101"; a(2)(6) <= "11001001"; a(2)(7) <= "11011101";

		b(2)(0) <= "10000110"; b(2)(1) <= "10001101"; b(2)(2) <= "11111111"; b(2)(3) <= "01110101";
		b(2)(4) <= "00000111"; b(2)(5) <= "10001101"; b(2)(6) <= "00001000"; b(2)(7) <= "11011110"; 
		
		-- ========================================================================================
		
		a(3)(0) <= "00000000"; a(3)(1) <= "00101010"; a(3)(2) <= "00000110"; a(3)(3) <= "00000000";
		a(3)(4) <= "00111010"; a(3)(5) <= "00001110"; a(3)(6) <= "00000000"; a(3)(7) <= "00000000";

		b(3)(0) <= "11111111"; b(3)(1) <= "00000111"; b(3)(2) <= "00000101"; b(3)(3) <= "00000001";
		b(3)(4) <= "00000111"; b(3)(5) <= "00000100"; b(3)(6) <= "00001000"; b(3)(7) <= "00000100"; 
		
		wait for 2 ns;
		end loop;
		wait;
	end process;
end;
