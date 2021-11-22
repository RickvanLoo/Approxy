library IEEE;
use IEEE.STD_LOGIC_1164.ALL;

entity uAcc2bitMult is
generic (word_size: integer:=2); 
Port ( 
A : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
B : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
prod: out STD_LOGIC_VECTOR (word_size * 2 - 1 downto 0));
end uAcc2bitMult;

architecture Behavioral of uAcc2bitMult is
begin
	process(A,B) is
	begin
		case A is
			when "00" =>
				case B is
                    when "00" =>
                        prod <= "0000";
                    when "01" =>
                        prod <= "0000";
                    when "10" =>
                        prod <= "0000";
                    when "11" =>
                        prod <= "0000";
                end case;
			when "01" =>
				case B is
                    when "00" =>
                        prod <= "0000";
                    when "01" =>
                        prod <= "0001";
                    when "10" =>
                        prod <= "0010";
                    when "11" =>
                        prod <= "0011";
                end case;
			when "10" =>
				case B is
                    when "00" =>
                        prod <= "0000";
                    when "01" =>
                        prod <= "0010";
                    when "10" =>
                        prod <= "0100";
                    when "11" =>
                        prod <= "0110";
                end case;
			when "11" =>
				case B is
                    when "00" =>
                        prod <= "0000";
                    when "01" =>
                        prod <= "0011";
                    when "10" =>
                        prod <= "0110";
                    when "11" =>
                        prod <= "1001";
                end case;
		end case;
		end process;

end Behavioral;