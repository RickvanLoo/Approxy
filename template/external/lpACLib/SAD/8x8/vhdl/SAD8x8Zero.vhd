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
-- SAD8x8Zero
-- Author: Walaa El-Harouni  
-- SAD8x8 built using AdderIMPACTZeroApproxMultiBit
-- Note: for changing the number of approx LSBs, edit this file and re-compile
-- uses: SAD8x1Zero, AdderIMPACTZeroApproxMultiBit, AdderIMPACTZeroApproxOneBit.vhd and AdderAccurateOneBit.vhd
--------------------------------------------

library ieee;
use ieee.std_logic_1164.ALL;
use ieee.numeric_std.ALL;
-- Define an array of 8 bytes
package TypesDefinition is
   type BYTE_ARRAY_8 is array(7 downto 0) of std_logic_vector(7 downto 0);
end TypesDefinition; 

library ieee;
use ieee.std_logic_1164.ALL;
use ieee.numeric_std.ALL;
use work.TypesDefinition.ALL;

entity SAD8x8Zero is
	generic (bitWidth : integer := 22; approxBits : integer := 6);
	port (A   	: in BYTE_ARRAY_8; --std_logic_vector(63 downto 0);	--! 8 bits per input
	      B 	: in BYTE_ARRAY_8;
	      i_valid	: in std_logic;
	      clk	: in std_logic;
	      reset	: in std_logic;
	      outReady  : out std_logic;
	      SAD 	: out std_logic_vector(bitWidth-1 downto 0));
end SAD8x8Zero;

architecture SAD8x8ZeroArch of SAD8x8Zero is
	component SAD8x1Zero is
		generic (bitWidth : integer := 12; approxBits : integer := 6);
		port (A   	 : in BYTE_ARRAY_8;
		      B 	 : in BYTE_ARRAY_8;
		      SAD 	 : out std_logic_vector(bitWidth-1 downto 0));
	end component SAD8x1Zero;

	component AdderIMPACTZeroApproxMultiBit is
		generic (bitWidth : integer := 22; approxBits : integer := 6);
		port (A   	 : in std_logic_vector(bitWidth-1 downto 0);
		      B 	 : in std_logic_vector(bitWidth-1 downto 0);
		      Cin 	 : in std_logic;
		      Sub 	 : in std_logic; -- '0' to add and '1' to subtract
		      Sum 	 : out std_logic_vector(bitWidth-1 downto 0);
		      Cout    : out std_logic );
	end component AdderIMPACTZeroApproxMultiBit;
	
	signal sadCycleOut	: std_logic_vector(11 downto 0);
	signal w_accumulatedSAD	: std_logic_vector(21 downto 0);
	signal w_final_addr_out	: std_logic_vector(21 downto 0);
	signal w_8x1_output	: std_logic_vector(21 downto 0);
	signal accumulatedSAD	: integer range 0 to 1048576; -- std_logic_vector(20 downto 0); -- equal to register size
	signal r_8x1_output	: integer range 0 to 1048576; -- std_logic_vector(sadCycleOut'range);
	signal r_valid 		: std_logic;

begin
	sad8x1Instance: SAD8x1Zero port map (A=>A, B=>B, SAD=>sadCycleOut);
	
	w_accumulatedSAd <= std_logic_vector(to_unsigned(accumulatedSAD, w_accumulatedSAD'length ));
	w_8x1_output <= std_logic_vector(to_unsigned(r_8x1_output, w_8x1_output'length ));
	
	FinalSADAcc : AdderIMPACTZeroApproxMultiBit
		generic map (approxBits => 6) 
		port map (A => w_accumulatedSAD,
			  B => w_8x1_output,
			  Cin => '0',
			  Sub => '0',
			  Sum => w_final_addr_out,
			  Cout => open);
	
	SAD <= std_logic_vector(to_unsigned(accumulatedSAD,SAD'length ));
	
	process(clk, reset)
		variable v_count : integer := 0;
	begin
		if reset = '1' then
			v_count := 0;
			 accumulatedSAD <= 0;	
			 outReady <= '0';	
			 r_valid <= '0';
		elsif (clk='1' and clk'event ) then
			r_valid <= i_valid;
			outReady <= '0';
			if i_valid = '1' then
				r_8x1_output <= to_integer(unsigned(sadCycleOut));
				if(v_count = 0) then
					accumulatedSAD <= 0;
				end if;
			end if;
			if r_valid = '1' then
				accumulatedSAD <= to_integer(unsigned(w_final_addr_out));
				v_count := v_count + 1;
				if(v_count = 8) then
					v_count := 0;
					outReady <= '1';
				end if;
			end if;			
		end if;
	end process;

end SAD8x8ZeroArch;




