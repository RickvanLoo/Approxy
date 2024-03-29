library IEEE;
use IEEE.STD_LOGIC_1164.ALL;
use ieee.numeric_std.all;

entity {{.EntityName}} is
generic (word_size: integer:=4); 
Port ( 
A : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
B : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
prod: out STD_LOGIC_VECTOR (word_size * 2 - 1 downto 0));
end {{.EntityName}};
    
architecture Behavioral of {{.EntityName}} is
    signal AH_BH_OUT, AH_BL_OUT, AL_BH_OUT, AL_BL_OUT :  std_logic_vector(3 downto 0) := (others => '0');
    signal AH, AL, BH, BL : std_logic_vector(1 downto 0) := (others => '0');
    signal factor1, factor2, factor3: unsigned(7 downto 0) := (others => '0');
begin

AH_BH: entity work.{{ (index .LUTArray 0).EntityName }} port map(a=>AH, b=>BH,prod=>AH_BH_OUT);
AH_BL: entity work.{{ (index .LUTArray 1).EntityName }} port map(a=>AH, b=>BL,prod=>AH_BL_OUT);
AL_BH: entity work.{{ (index .LUTArray 2).EntityName }} port map(a=>AL, b=>BH,prod=>AL_BH_OUT);
AL_BL: entity work.{{ (index .LUTArray 3).EntityName }} port map(a=>AL, b=>BL,prod=>AL_BL_OUT);

-- //LUTArray[0] = AH*BH
-- //LUTArray[1] = AH*BL
-- //LUTArray[2] = AL*BH
-- //LUTArray[3] = AL*BL

AH <= A(3 downto 2);
AL <= A(1 downto 0);
BH <= B(3 downto 2);
BL <= B(1 downto 0);

-- factor1 <= unsigned(AL_BL_OUT);
-- factor2 <= shift_left(unsigned(AL_BH_OUT),2);
-- factor3 <= shift_left(unsigned(AH_BL_OUT),2);
-- factor4 <= shift_left(unsigned(AH_BH_OUT),4);

-- factor1 <= "0000" & unsigned(AL_BL_OUT);
factor2 <= "00" & unsigned(AL_BH_OUT) & "00";
factor3 <= "00" & unsigned(AH_BL_OUT) & "00";
-- factor4 <= unsigned(AH_BH_OUT) & "0000";
factor1 <= unsigned(AH_BH_OUT) & unsigned(AL_BL_OUT); 

prod <= std_logic_vector(factor1 + factor2 + factor3);

end Behavioral;