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
-- SAD32x32Zero
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
   type BYTE_ARRAY_32 is array(3 downto 0) of BYTE_ARRAY_8;
end TypesDefinition; 

library ieee;
use ieee.std_logic_1164.ALL;
use ieee.numeric_std.ALL;
use work.TypesDefinition.ALL;

entity SAD32x32Zero is
	generic (bitWidth : integer := 22; approxBits : integer := 6);
	port (A32   	: in BYTE_ARRAY_32;
	      B32 	: in BYTE_ARRAY_32;
	      i_valid	: in std_logic;
	      clk	: in std_logic;
	      reset	: in std_logic;
	      outReady  : out std_logic;
	      SAD 	: out std_logic_vector(bitWidth-1 downto 0));
end SAD32x32Zero;

architecture SAD32x32ZeroArch of SAD32x32Zero is
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
	
	signal sadCycleOut1	: std_logic_vector(11 downto 0);
	signal sadCycleOut2	: std_logic_vector(11 downto 0);
	signal sadCycleOut3	: std_logic_vector(11 downto 0);
	signal sadCycleOut4	: std_logic_vector(11 downto 0);
	
	signal sadResized1	: std_logic_vector(21 downto 0);
	signal sadResized2	: std_logic_vector(21 downto 0);
	signal sadResized3	: std_logic_vector(21 downto 0);
	signal sadResized4	: std_logic_vector(21 downto 0);
	
	signal L1Add0		: std_logic_vector(21 downto 0);
	signal L1Add1		: std_logic_vector(21 downto 0);
	signal L2Add0		: std_logic_vector(21 downto 0);
	
	signal w_accumulatedSAD	: std_logic_vector(21 downto 0);
	signal w_final_addr_out	: std_logic_vector(21 downto 0);
	signal w_8x1_output	: std_logic_vector(21 downto 0);
	signal accumulatedSAD	: integer range 0 to 1048576; -- std_logic_vector(20 downto 0); -- equal to register size
	signal r_8x1_output	: integer range 0 to 1048576; -- std_logic_vector(sadCycleOut'range);
	signal r_valid 		: std_logic;

begin
	sad8x1Instance1: SAD8x1Zero port map (A=>A32(0), B=>B32(0), SAD => sadCycleOut1);
	sad8x1Instance2: SAD8x1Zero port map (A=>A32(1), B=>B32(1), SAD => sadCycleOut2);
	sad8x1Instance3: SAD8x1Zero port map (A=>A32(2), B=>B32(2), SAD => sadCycleOut3);
	sad8x1Instance4: SAD8x1Zero port map (A=>A32(3), B=>B32(3), SAD => sadCycleOut4);
	
	sadResized1 <= std_logic_vector( resize( unsigned(sadCycleOut1), sadResized1'length) );
	sadResized2 <= std_logic_vector( resize( unsigned(sadCycleOut2), sadResized2'length) );
	sadResized3 <= std_logic_vector( resize( unsigned(sadCycleOut3), sadResized3'length) );
	sadResized4 <= std_logic_vector( resize( unsigned(sadCycleOut4), sadResized4'length) );

	Level1_Adder0 : AdderIMPACTZeroApproxMultiBit
		generic map (approxBits => 6) 
		port map (A => sadResized1,
			  B => sadResized2,
			      Cin => '0',
			      Sub => '0',
			      Sum => L1Add0,
			      Cout => open);
			      
	Level1_Adder1 : AdderIMPACTZeroApproxMultiBit
		generic map (approxBits => 6) 
		port map (A => sadResized3,
			  B => sadResized4,
			      Cin => '0',
			      Sub => '0',
			      Sum => L1Add1,
			      Cout => open);
	
	Level2_Adder0 : AdderIMPACTZeroApproxMultiBit
		generic map (approxBits => 6) 
		port map (A => L1Add0,
			  B => L1Add1,
			      Cin => '0',
			      Sub => '0',
			      Sum => L2Add0,
			      Cout => open);
	
	
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
				r_8x1_output <= to_integer(unsigned(L2Add0));
				if(v_count = 0) then
					accumulatedSAD <= 0;
				end if;
			end if;
			if r_valid = '1' then
				accumulatedSAD <= to_integer(unsigned(w_final_addr_out));
				v_count := v_count + 1;
				if(v_count = 32) then
					v_count := 0;
					outReady <= '1';
				end if;
			end if;			
		end if;
	end process;

end SAD32x32ZeroArch;




