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
-- SAD8x1First
-- Author: Walaa El-Harouni  
-- SAD8x1 built using AdderIMPACTFirstApproxMultiBit
-- Note: for changing the number of approx LSBs, edit this file and re-compile
-- uses: AdderIMPACTFirstApproxMultiBit, AdderIMPACTFirstApproxOneBit.vhd and AdderAccurateOneBit.vhd
--------------------------------------------

library ieee;
use ieee.std_logic_1164.ALL;
use ieee.numeric_std.ALL;
package TypesDefinition is
   type BYTE_ARRAY_8 is array(7 downto 0) of std_logic_vector(7 downto 0);
   type BYTE_ARRAY_32 is array(3 downto 0) of BYTE_ARRAY_8;
end TypesDefinition; 

library ieee;
use ieee.std_logic_1164.ALL;
use ieee.numeric_std.ALL;
use work.TypesDefinition.ALL;

entity SAD8x1First is
	generic (bitWidth : integer := 12; approxBits : integer := 6);
	port (A   	 : in BYTE_ARRAY_8;
	      B 	 : in BYTE_ARRAY_8;
	     SAD 	 : out std_logic_vector(bitWidth-1 downto 0));
end SAD8x1First;

architecture SAD8x1FirstArch of SAD8x1First is
	component AdderIMPACTFirstApproxMultiBit
		generic (bitWidth : integer := 9; approxBits : integer := 6);
		port (A   	 : in std_logic_vector(bitWidth-1 downto 0);
			  B 	 : in std_logic_vector(bitWidth-1 downto 0);
			  Cin 	 : in std_logic;
			  Sub 	 : in std_logic; -- '0' to add and '1' to subtract
			  Sum 	 : out std_logic_vector(bitWidth-1 downto 0);
			  Cout   : out std_logic );
	end component;

 	-- Level1 signals (8 adders)
	 type TYPE_SUM_L1 is array (0 to 7) of std_logic_vector(8 downto 0); -- 9-bit output
	 signal Sum_L1: TYPE_SUM_L1; 
	 signal Cout_L1: std_logic_vector(0 to 7);
	 
	 type TYPE_ABS is array (0 to 7) of std_logic_vector(9 downto 0);
	 signal absoluteDiff: TYPE_ABS;
	  
	 
	-- Level2 signals (4 adders)
	 type TYPE_SUM_L2 is array (0 to 3) of std_logic_vector(9 downto 0); -- 10-bit output
	 signal Sum_L2: TYPE_SUM_L2; 
	 signal Cout_L2: std_logic_vector(0 to 3);
	 
	 -- Level3 signals (2 adders)
	 type TYPE_SUM_L3 is array (0 to 1) of std_logic_vector(10 downto 0); -- 11-bit output
	 signal Sum_L3: TYPE_SUM_L3; 
	 signal Cout_L3: std_logic_vector(0 to 1);
	 
	 -- Level4 signals (1 adder)
	 signal Sum_L4: std_logic_vector(11 downto 0);  --12-bit output
	 signal Cout_L4: std_logic;
	 
	 -- for resizing
	 signal L1_A: TYPE_SUM_L1; signal L1_B: TYPE_SUM_L1;
	 signal L3_A: TYPE_SUM_L3; signal L3_B: TYPE_SUM_L3;
	 signal L4_A: std_logic_vector(11 downto 0); signal L4_B: std_logic_vector(11 downto 0);
	 
begin

  -- resize inputs for level 1
      L1_resize : process(A, B)
      begin
	      for i in A'range loop
		      L1_A(i) <= std_logic_vector( resize( signed(A(i)), 9) );
		      L1_B(i) <= std_logic_vector( resize( signed(B(i)), 9) );
	      end loop;
      end process L1_resize;
	
-- Level 1 => 8 adders (actually subtractors)
  levelOne: for i in 0 to 7 generate 
    levelOneInstances: AdderIMPACTFirstApproxMultiBit
      generic map(bitWidth => 9)
      port map (
        A  => L1_A(i),
        B  => L1_B(i),
	Sum => Sum_L1(i), 
        Cin => '0',
        Sub => '1',
        Cout => Cout_L1(i)
      );
  end generate;
  
 -- calculating absolutes (resize in place)
 	absolutes : process(Sum_L1)
	begin
		for i in Sum_L1'range loop
			absoluteDiff(i) <= std_logic_vector( resize(abs(signed(Sum_L1(i))), 10) );
		end loop;
	end process absolutes;
  
  --  Level 2 => 4 adders
  -- Adder0: Sum_L1(0) and Sum_L1(1)
  -- Adder1: Sum_L1(2) and Sum_L1(3)
  -- Adder2: Sum_L1(4) and Sum_L1(5)
  -- Adder3: Sum_L1(6) and Sum_L1(7)
  levelTwo: for i in 0 to 3 generate 
    levelTwoInstances: AdderIMPACTFirstApproxMultiBit
      generic map(bitWidth => 10)
      port map (
        A  => absoluteDiff(i+i*1),
        B  => absoluteDiff(i+1+i*1),
	Sum => Sum_L2(i), 
	Cin => '0',
	Sub => '0',
	Cout => Cout_L2(i)
      );
  end generate;
  
  
    -- resize inputs for level 3
     L3_A(0) <= std_logic_vector( resize( signed(Sum_L2(0)), 11) );
     L3_B(0) <= std_logic_vector( resize( signed(Sum_L2(1)), 11) );
     L3_A(1) <= std_logic_vector( resize( signed(Sum_L2(2)), 11) );
     L3_B(1) <= std_logic_vector( resize( signed(Sum_L2(3)), 11) );
      
  -- Level 3 => 2 adders  
 levelThreeAdder0: AdderIMPACTFirstApproxMultiBit generic map(bitWidth => 11)
 	port map (A  => L3_A(0),
		  B  => L3_B(0), 
		  Sum => Sum_L3(0), 
		  Cin => '0', 
		  Sub => '0', 
		  Cout => Cout_L3(0));
  
 levelThreeAdder2: AdderIMPACTFirstApproxMultiBit generic map(bitWidth => 11)
	 port map (A  => L3_A(1),
		   B  => L3_B(1),
		   Sum => Sum_L3(1),
		   Cin => '0', 
		   Sub => '0', 
		   Cout => Cout_L3(1));
  
  
  -- resize inputs for level 4
  L4_A <= std_logic_vector( resize(signed(Sum_L3(0)), 12));
  L4_B <= std_logic_vector( resize(signed(Sum_L3(1)), 12));
  
 -- Level 4 => 1 adders 
 levelFourAdder: AdderIMPACTFirstApproxMultiBit generic map(bitWidth => 12)
      port map (A  => L4_A,
		B  => L4_B,
		Sum => Sum_L4, 
		Cin => '0', 
		Sub => '0', 
		Cout => Cout_L4);
        
      SAD <= Sum_L4;
	
end SAD8x1FirstArch;




